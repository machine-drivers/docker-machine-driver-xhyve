package xhyve

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/mcnutils"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	"github.com/johanneswuerbach/nfsexports"
	"github.com/zchee/docker-machine-driver-xhyve/vmnet"
)

const (
	isoFilename           = "boot2docker.iso"
	isoMountPath          = "b2d-image"
	defaultBoot2DockerURL = ""
	defaultBootCmd        = "loglevel=3 user=docker console=ttyS0 console=tty0 noembed nomodeset norestore waitusb=10 base host=boot2docker"
	defaultCPU            = 1
	defaultCaCertPath     = ""
	defaultDiskSize       = 20000
	defaultMacAddr        = ""
	defaultMemory         = 1024
	defaultPrivateKeyPath = ""
	defaultUUID           = ""
	defaultNFSShare       = false
)

type Driver struct {
	*drivers.BaseDriver
	Boot2DockerURL string
	BootCmd        string
	CPU            int
	CaCertPath     string
	DiskSize       int64
	MacAddr        string
	Memory         int
	PrivateKeyPath string
	UUID           string
	NFSShare       bool
}

var (
	ErrMachineExist    = errors.New("machine already exists")
	ErrMachineNotExist = errors.New("machine does not exist")
)

// NewDriver creates a new VirtualBox driver with default settings.
func NewDriver(hostName, storePath string) *Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
		Boot2DockerURL: defaultBoot2DockerURL,
		BootCmd:        defaultBootCmd,
		CPU:            defaultCPU,
		CaCertPath:     defaultCaCertPath,
		DiskSize:       defaultDiskSize,
		MacAddr:        defaultMacAddr,
		Memory:         defaultMemory,
		PrivateKeyPath: defaultPrivateKeyPath,
		UUID:           defaultUUID,
		NFSShare:       defaultNFSShare,
	}
}

// RegisterCreateFlags registers the flags this driver adds to
// "docker hosts create"
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "XHYVE_BOOT2DOCKER_URL",
			Name:   "xhyve-boot2docker-url",
			Usage:  "The URL of the boot2docker image. Defaults to the latest available version",
			Value:  "",
		},
		mcnflag.IntFlag{
			EnvVar: "XHYVE_CPU_COUNT",
			Name:   "xhyve-cpu-count",
			Usage:  "Number of CPUs for the machine (-1 to use the number of CPUs available)",
			Value:  defaultCPU,
		},
		mcnflag.IntFlag{
			EnvVar: "XHYVE_MEMORY_SIZE",
			Name:   "xhyve-memory-size",
			Usage:  "Size of memory for host in MB",
			Value:  defaultMemory,
		},
		mcnflag.IntFlag{
			EnvVar: "XHYVE_DISK_SIZE",
			Name:   "xhyve-disk-size",
			Usage:  "Size of disk for host in MB",
			Value:  defaultDiskSize,
		},
		mcnflag.StringFlag{
			EnvVar: "XHYVE_BOOT_CMD",
			Name:   "xhyve-boot-cmd",
			Usage:  "Command of booting kexec protocol",
			Value:  defaultBootCmd,
		},
		mcnflag.BoolFlag{
			EnvVar: "XHYVE_EXPERIMENTAL_NFS_SHARE",
			Name:   "xhyve-experimental-nfs-share",
			Usage:  "Setup NFS shared folder (requires root)",
		},
	}
}

func (d *Driver) GetMachineName() string {
	return d.MachineName
}

func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

func (d *Driver) GetSSHKeyPath() string {
	return d.ResolveStorePath("id_rsa")
}

func (d *Driver) GetSSHPort() (int, error) {
	if d.SSHPort == 0 {
		d.SSHPort = 22
	}

	return d.SSHPort, nil
}

func (d *Driver) GetSSHUsername() string {
	if d.SSHUser == "" {
		d.SSHUser = "docker"
	}

	return d.SSHUser
}

func (d *Driver) DriverName() string {
	return "xhyve"
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.Boot2DockerURL = flags.String("xhyve-boot2docker-url")
	d.CPU = flags.Int("xhyve-cpu-count")
	d.Memory = flags.Int("xhyve-memory-size")
	d.DiskSize = int64(flags.Int("xhyve-disk-size"))
	d.BootCmd = flags.String("xhyve-boot-cmd")
	d.SwarmMaster = flags.Bool("swarm-master")
	d.SwarmHost = flags.String("swarm-host")
	d.SwarmDiscovery = flags.String("swarm-discovery")
	d.SSHUser = "docker"
	d.SSHPort = 22
	d.NFSShare = flags.Bool("xhyve-experimental-nfs-share")

	return nil
}

func (d *Driver) GetURL() (string, error) {
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	if ip == "" {
		return "", nil
	}
	return fmt.Sprintf("tcp://%s:2376", ip), nil
}

func (d *Driver) GetIP() (string, error) {
	if d.IPAddress != "" {
		return d.IPAddress, nil
	}

	return d.getIPfromDHCPLease()
}

func (d *Driver) GetState() (state.State, error) {
	s, _ := d.GetSShState()
	if !s {
		return state.Stopped, nil
	}
	return state.Running, nil
}

func (d *Driver) GetSShState() (bool, error) {
	log.Debug("Getting to VM SSH status...")
	if _, err := drivers.RunSSHCommandFromDriver(d, "exit 0"); err != nil {
		return false, nil
	}
	return true, nil
}

// Print driver version, Check VirtualBox version
func (d *Driver) PreCreateCheck() error {
	//TODO:libmachine PLEASE output driver version API!
	v := Version
	c := GitCommit
	log.Debugf("===== Docker Machine %s Driver Version %s (%s) =====\n", d.DriverName(), v, c)

	ver, err := vboxVersionDetect()
	if ver == "" && err == nil {
		return nil
	}
	if err != nil {
		return fmt.Errorf("Error detecting VBox version: %s", err)
	}
	if !strings.HasPrefix(ver, "5") {
		return fmt.Errorf("Virtual Box version 4 or lower will cause a kernel panic" +
			"if xhyve tries to run. You are running version: " +
			ver +
			"\n\t Please upgrade to version 5 at https://www.virtualbox.org/wiki/Downloads")
	}
	return nil
}

func (d *Driver) Create() error {
	b2dutils := mcnutils.NewB2dUtils(d.StorePath)
	if err := b2dutils.CopyIsoToMachineDir(d.Boot2DockerURL, d.MachineName); err != nil {
		return err
	}

	log.Infof("Creating SSH key...")
	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return err
	}

	log.Infof("Creating VM...")
	if err := os.MkdirAll(d.ResolveStorePath("."), 0755); err != nil {
		return err
	}

	log.Infof("Extracting vmlinuz64 and initrd.img from %s...", isoFilename)
	if err := d.extractKernelImages(); err != nil {
		return err
	}

	log.Infof("Generating %dMB disk image...", d.DiskSize)
	if err := d.generateDiskImage(d.DiskSize); err != nil {
		return err
	}

	// Fix file permission root to current user for vmnet.framework
	log.Infof("Fix file permission...")
	os.Chown(d.ResolveStorePath("."), syscall.Getuid(), syscall.Getegid())
	files, _ := ioutil.ReadDir(d.ResolveStorePath("."))
	for _, f := range files {
		log.Debugf(d.ResolveStorePath(f.Name()))
		os.Chown(d.ResolveStorePath(f.Name()), syscall.Getuid(), syscall.Getegid())
	}

	log.Infof("Generate UUID...")
	d.UUID = uuidgen()
	log.Debugf("Generated UUID: %s", d.UUID)

	log.Infof("Convert UUID to MAC address...")
	rawUUID, err := d.getMACAdress()
	if err != nil {
		return err
	}
	d.MacAddr = trimMacAddress(rawUUID)
	log.Debugf("Converted MAC address: %s", d.MacAddr)

	log.Infof("Starting %s...", d.MachineName)
	if err := d.Start(); err != nil {
		return err
	}
	log.Infof("Waiting for VM to come online...")

	var ip string
	for i := 1; i <= 60; i++ {
		ip, err = d.getIPfromDHCPLease()
		if err != nil {
			log.Debugf("Not there yet %d/%d, error: %s", i, 60, err)
			time.Sleep(2 * time.Second)
			continue
		}

		if ip != "" {
			log.Debugf("Got an ip: %s", ip)
			break
		}
	}

	if ip == "" {
		return fmt.Errorf("Machine didn't return an IP after 120 seconds, aborting")
	}

	// We got an IP, let's copy ssh keys over
	d.IPAddress = ip

	// Setup NFS sharing
	if d.NFSShare {
		err = d.setupNFSShare()
		if err != nil {
			log.Errorf("NFS setup failed: %s", err.Error())
		}
	}

	return nil
}

func (d *Driver) Start() error {
	args := d.xhyveArgs()

	log.Debug(args)

	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	err := cmd.Start()
	if err != nil {
		return err
	}

	go func() {
		err := cmd.Wait()
		if err != nil {
			log.Error(err, cmd.Stdout, cmd.Stderr)
		}
	}()

	return nil
}

func (d *Driver) Stop() error {
	log.Infof("Stopping %s use send ACPI signals poweroff ...", d.MachineName)
	if _, err := drivers.RunSSHCommandFromDriver(d, "sudo poweroff"); err != nil {
		log.Debugf("Error getting ssh command 'exit 0' : %s", err)
		return err
	}

	for {
		s, err := d.GetState()
		if err != nil {
			return err
		}
		if s == state.Running {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	d.IPAddress = ""

	return nil
}

func (d *Driver) Remove() error {
	s, err := d.GetState()
	if err != nil {
		if err == ErrMachineNotExist {
			log.Infof("machine does not exist, assuming it has been removed already")
			return nil
		}
		return err
	}
	if s == state.Running {
		if err := d.Stop(); err != nil {
			return err
		}

		if _, err := nfsexports.Remove("", d.nfsExportIdentifier()); err != nil {
			log.Errorf("failed removing nfs share: %s", err.Error())
		}

		if err := nfsexports.ReloadDaemon(); err != nil {
			log.Errorf("failed reload nfs daemon: %s", err.Error())
		}
	}
	return nil
}

func (d *Driver) Restart() error {
	s, err := d.GetState()
	if err != nil {
		return err
	}

	if s == state.Running {
		if err := d.Stop(); err != nil {
			return err
		}
	}
	return d.Start()
}

func (d *Driver) Kill() error {
	log.Infof("Killing %s use hardware to stop all CPU ...", d.MachineName)
	if _, err := drivers.RunSSHCommandFromDriver(d, "sudo halt"); err != nil {
		log.Debugf("Error getting ssh command 'exit 0' : %s", err)
		return err
	}

	return nil
}

func (d *Driver) setMachineNameIfNotSet() {
	if d.MachineName == "" {
		d.MachineName = fmt.Sprintf("docker-machine-unknown")
	}
}

func (d *Driver) getIPfromDHCPLease() (string, error) {
	currentip, err := vmnet.GetIPAddressByMACAddress(d.MacAddr)
	log.Debugf(currentip)

	if currentip == "" {
		return "", fmt.Errorf("IP not found for MAC %s in DHCP leases", d.MacAddr)
	}

	log.Debugf("IP found in DHCP lease table: %s", currentip)
	return currentip, err
}

func (d *Driver) publicSSHKeyPath() string {
	return d.GetSSHKeyPath() + ".pub"
}

func (d *Driver) extractKernelImages() error {
	log.Debugf("Mounting %s", isoFilename)

	err := hdiutil("attach", d.ResolveStorePath(isoFilename), "-mountpoint", d.ResolveStorePath("b2d-image"))
	if err != nil {
		return err
	}

	volumeRootDir := d.ResolveStorePath(isoMountPath)
	vmlinuz64 := volumeRootDir + "/boot/vmlinuz64"
	initrd := volumeRootDir + "/boot/initrd.img"

	log.Debugf("Extracting vmlinuz64 into %s", d.ResolveStorePath("."))
	if err := mcnutils.CopyFile(vmlinuz64, d.ResolveStorePath("vmlinuz64")); err != nil {
		return err
	}
	log.Debugf("Extracting initrd.img into %s", d.ResolveStorePath("."))
	if err := mcnutils.CopyFile(initrd, d.ResolveStorePath("initrd.img")); err != nil {
		return err
	}
	log.Debugf("Unmounting %s", isoFilename)
	if err := hdiutil("detach", volumeRootDir); err != nil {
		return err
	}

	return nil
}

func (d *Driver) generateDiskImage(count int64) error {
	output := d.ResolveStorePath(d.MachineName)

	if err := hdiutil("create", "-megabytes", fmt.Sprintf("%d", d.DiskSize), output); err != nil {
		return err
	}

	tarBuf, err := d.generateKeyBundle()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(output+".dmg", os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Seek(0, os.SEEK_SET)
	_, err = file.Write(tarBuf.Bytes())
	if err != nil {
		return err
	}
	file.Close()

	return nil
}

// Make a boot2docker userdata.tar key bundle
func (d *Driver) generateKeyBundle() (*bytes.Buffer, error) {
	magicString := "boot2docker, please format-me"

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// magicString first so the automount script knows to format the disk
	file := &tar.Header{Name: magicString, Size: int64(len(magicString))}
	if err := tw.WriteHeader(file); err != nil {
		return nil, err
	}
	if _, err := tw.Write([]byte(magicString)); err != nil {
		return nil, err
	}
	// .ssh/key.pub => authorized_keys
	file = &tar.Header{Name: ".ssh", Typeflag: tar.TypeDir, Mode: 0700}
	if err := tw.WriteHeader(file); err != nil {
		return nil, err
	}
	pubKey, err := ioutil.ReadFile(d.publicSSHKeyPath())
	if err != nil {
		return nil, err
	}
	file = &tar.Header{Name: ".ssh/authorized_keys", Size: int64(len(pubKey)), Mode: 0644}
	if err := tw.WriteHeader(file); err != nil {
		return nil, err
	}
	if _, err := tw.Write([]byte(pubKey)); err != nil {
		return nil, err
	}
	file = &tar.Header{Name: ".ssh/authorized_keys2", Size: int64(len(pubKey)), Mode: 0644}
	if err := tw.WriteHeader(file); err != nil {
		return nil, err
	}
	if _, err := tw.Write([]byte(pubKey)); err != nil {
		return nil, err
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	return buf, nil
}

// Setup NFS share
func (d *Driver) setupNFSShare() error {
	nfsConfig := fmt.Sprintf("/Users %s -alldirs -maproot=root", d.IPAddress)

	if _, err := nfsexports.Add("", d.nfsExportIdentifier(), nfsConfig); err != nil {
		return err
	}

	if err := nfsexports.ReloadDaemon(); err != nil {
		return err
	}

	hostIP, err := vmnet.GetNetAddr()
	if err != nil {
		return err
	}

	bootScriptName := "/var/lib/boot2docker/bootlocal.sh"
	bootScript := fmt.Sprintf("#/bin/bash\\n"+
		"sudo mkdir -p /Users\\n"+
		"sudo /usr/local/etc/init.d/nfs-client start\\n"+
		"sudo mount -t nfs -o noacl,async %s:/Users /Users\\n", hostIP)

	writeScriptCmd := fmt.Sprintf("echo -e \"%s\" | sudo tee %s && sudo chmod +x %s && %s",
		bootScript, bootScriptName, bootScriptName, bootScriptName)

	if _, err := drivers.RunSSHCommandFromDriver(d, writeScriptCmd); err != nil {
		return err
	}

	return nil
}

func (d *Driver) nfsExportIdentifier() string {
	return fmt.Sprintf("docker-machine-driver-xhyve %s", d.MachineName)
}

//Trimming "0" of the ten's digit
func trimMacAddress(rawUUID string) string {
	re := regexp.MustCompile(`[0]([A-Fa-f0-9][:])`)
	mac := re.ReplaceAllString(rawUUID, "$1")

	return mac
}

func (d *Driver) xhyveArgs() []string {
	uuid := d.UUID
	vmlinuz := d.ResolveStorePath("vmlinuz64")
	initrd := d.ResolveStorePath("initrd.img")
	iso := d.ResolveStorePath(isoFilename)
	img := d.ResolveStorePath(d.MachineName + ".dmg")
	bootcmd := d.BootCmd

	cpus := d.CPU
	if cpus < 1 {
		cpus = int(runtime.NumCPU())
	}

	return []string{
		"xhyve",
		"-A",
		"-U", fmt.Sprintf("%s", uuid),
		"-c", fmt.Sprintf("%d", cpus),
		"-m", fmt.Sprintf("%dM", d.Memory),
		"-l", "com1",
		"-s", "0:0,hostbridge",
		"-s", "31,lpc",
		"-s", "2:0,virtio-net",
		"-s", fmt.Sprintf("3,ahci-cd,%s", iso),
		"-s", fmt.Sprintf("4,virtio-blk,%s", img),
		"-f", fmt.Sprintf("kexec,%s,%s,%s", vmlinuz, initrd, bootcmd)}
}

func (d *Driver) getMACAdress() (string, error) {
	args := append(d.xhyveArgs(), "-M")

	stdout := bytes.Buffer{}

	cmd := exec.Command(os.Args[0], args...) // TODO: Should be possible without exec
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	mac := bytes.TrimPrefix(stdout.Bytes(), []byte("MAC: "))
	mac = bytes.TrimSpace(mac)

	hw, err := net.ParseMAC(string(mac))
	if err != nil {
		return "", err
	}
	return hw.String(), nil
}
