package xhyve

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/mcnutils"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	"github.com/zchee/docker-machine-xhyve/version"
	"github.com/zchee/docker-machine-xhyve/vmnet"
	"libguestfs.org/guestfs"
)

const (
	isoFilename = "boot2docker.iso"
)

type Driver struct {
	*drivers.BaseDriver
	Memory         int
	DiskSize       int64
	CPU            int
	TmpISO         string
	UUID           string
	MacAddr        string
	BootCmd        string
	Boot2DockerURL string
	CaCertPath     string
	PrivateKeyPath string
}

var (
	ErrMachineExist    = errors.New("machine already exists")
	ErrMachineNotExist = errors.New("machine does not exist")
)

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
			Value:  1,
		},
		mcnflag.IntFlag{
			EnvVar: "XHYVE_MEMORY_SIZE",
			Name:   "xhyve-memory",
			Usage:  "Size of memory for host in MB",
			Value:  1024,
		},
		mcnflag.IntFlag{
			EnvVar: "XHYVE_DISK_SIZE",
			Name:   "xhyve-disk-size",
			Usage:  "Size of disk for host in MB",
			Value:  20000,
		},
		mcnflag.StringFlag{
			EnvVar: "XHYVE_BOOT_CMD",
			Name:   "xhyve-boot-cmd",
			Usage:  "Command of booting kexec protocol",
			Value:  "loglevel=3 user=docker console=ttyS0 console=tty0 noembed nomodeset norestore waitusb=10:LABEL=boot2docker-data base host=boot2docker",
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
	d.Memory = flags.Int("xhyve-memory")
	d.DiskSize = int64(flags.Int("xhyve-disk-size"))
	d.BootCmd = flags.String("xhyve-boot-cmd")
	d.SwarmMaster = flags.Bool("swarm-master")
	d.SwarmHost = flags.String("swarm-host")
	d.SwarmDiscovery = flags.String("swarm-discovery")
	d.SSHUser = "docker"
	d.SSHPort = 22

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
	s, err := d.GetState()
	if err != nil {
		return "", err
	}
	if s != state.Running {
		return "", drivers.ErrHostIsNotRunning
	}

	ip, err := d.getIPfromDHCPLease()
	if err != nil {
		return "", err
	}

	return ip, nil
}

func (d *Driver) GetState() (state.State, error) { // TODO
	// VMRUN only tells use if the vm is running or not
	//	if stdout, _, _ := vmrun("list"); strings.Contains(stdout, d.vmxPath()) {
	return state.Running, nil
	//	}
	// return state.Stopped, nil
}

// Check VirtualBox version
func (d *Driver) PreCreateCheck() error {
	//TODO:libmachine PLEASE output driver version API!
	v := version.Version
	c := version.GitCommit
	log.Debugf("===== Docker Machine %s Driver Version %s (%s) =====\n", d.DriverName(), v, c)

	ver, err := vboxVersionDetect()
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
	b2dutils := mcnutils.NewB2dUtils("", "", d.StorePath)
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

	log.Infof("Make a boot2docker userdata.tar key bundle...")
	if err := d.generateKeyBundle(); err != nil {
		return err
	}

	// Fix file permission root to current user.
	// In order to avoid require sudo of vmnet.framework, Execute the root owner(and root uid)
	// "docker-machine-xhyve" and "goxhyve" binary in golang.
	log.Infof("Fix file permission...")
	os.Chown(d.ResolveStorePath("."), 501, 20) //TODO Parse current user uid and gid
	files, _ := ioutil.ReadDir(d.ResolveStorePath("."))
	for _, f := range files {
		log.Debugf(d.ResolveStorePath(f.Name()))
		os.Chown(d.ResolveStorePath(f.Name()), 501, 20)
	}

	log.Infof("Creating blank ext4 filesystem disk image...")
	if err := d.generateBlankDiskImage(d.DiskSize); err != nil {
		return err
	}
	os.Chown(d.ResolveStorePath(d.MachineName+".img"), 501, 20)
	log.Debugf("Created disk size: %dMB", d.DiskSize)

	log.Infof("Generate UUID...")
	d.UUID = uuidgen() //TODO Native golang instead execute "uuidgen"
	log.Debugf("uuidgen generated UUID: %s", d.UUID)

	log.Infof("Convert UUID to MAC address...")
	//Trim "0" string
	re := regexp.MustCompile(`[0]`)
	mac, _ := vmnet.GetMACAddressByUUID(d.UUID)
	d.MacAddr = re.ReplaceAllString(mac, "")
	log.Debugf("uuid2mac output MAC address: %s", d.MacAddr)

	log.Infof("Starting %s...", d.MachineName)
	if err := d.Start(); err != nil {
		return err
	}

	var ip string
	var err error
	log.Infof("Waiting for VM to come online...")
	log.Debugf("d.MacAddr", d.MacAddr)
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

	// we got an IP, let's copy ssh keys over
	d.IPAddress = ip

	return nil
}

func (d *Driver) Start() error {
	uuid := d.UUID
	vmlinuz := d.ResolveStorePath("vmlinuz64")
	initrd := d.ResolveStorePath("initrd.img")
	iso := d.ResolveStorePath(isoFilename)
	img := d.ResolveStorePath(d.MachineName + ".img")
	bootcmd := d.BootCmd

	cmd := exec.Command("goxhyve",
		fmt.Sprintf("%s", uuid),
		fmt.Sprintf("%d", d.CPU),
		fmt.Sprintf("%d", d.Memory),
		fmt.Sprintf("%s", iso),
		fmt.Sprintf("%s", img),
		fmt.Sprintf("kexec,%s,%s,%s", vmlinuz, initrd, bootcmd),
		"-d", //TODO fix daemonize flag
	)
	log.Debug(cmd)
	go func() {
		err := cmd.Run()
		if err != nil {
			log.Error(err, cmd.Stdout)
		}
	}()

	return nil
}

func (d *Driver) Stop() error { // TODO
	// xhyve("controlvm", d.MachineName, "acpipowerbutton")
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

func (d *Driver) Remove() error { // TODO
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
	}
	//return xhyve("unregistervm", "--delete", d.MachineName)
	return nil
}

func (d *Driver) Restart() error { // TODO
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

func (d *Driver) Kill() error { // TODO
	//return xhyve("controlvm", d.MachineName, "poweroff")
	return nil
}

func (d *Driver) setMachineNameIfNotSet() {
	if d.MachineName == "" {
		d.MachineName = fmt.Sprintf("docker-machine-unknown")
	}
}

func (d *Driver) imgPath() string {
	return d.ResolveStorePath(fmt.Sprintf("%s.img", d.MachineName))
}

func (d *Driver) userdataPath() string {
	return d.ResolveStorePath("userdata.tar")
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
	var vmlinuz64 = "/Volumes/Boot2Docker-v1.8/boot/vmlinuz64" // TODO Do not hardcode boot2docker version
	var initrd = "/Volumes/Boot2Docker-v1.8/boot/initrd.img"   // TODO Do not hardcode boot2docker version

	log.Debugf("Mounting %s", isoFilename)
	hdiutil("attach", d.ResolveStorePath(isoFilename)) // TODO need parse attached disk identifier.

	log.Debugf("Extract vmlinuz64")
	if err := mcnutils.CopyFile(vmlinuz64, d.ResolveStorePath("vmlinuz64")); err != nil {
		return err
	}
	log.Debugf("Extract initrd.img")
	if err := mcnutils.CopyFile(initrd, d.ResolveStorePath("initrd.img")); err != nil {
		return err
	}
	log.Debugf("Unmounting %s", isoFilename)
	if err := hdiutil("detach", "/Volumes/Boot2Docker-v1.8/"); err != nil { // TODO Do not hardcode boot2docker version
		return err
	}

	return nil
}

func (d *Driver) generateBlankDiskImage(count int64) error {
	output := d.ResolveStorePath(d.MachineName + ".img")

	g, errno := guestfs.Create()
	if errno != nil {
		panic(errno)
	}
	defer g.Close()

	/* Set $LIBGUESTFS_PATH(libguestfs appliance path) to root user */
	p := toPtr("/usr/local/lib/guestfs")
	g.Set_path(p)

	/* Set the trace flag so that we can see each libguestfs call. */
	if log.IsDebug == true {
		g.Set_trace(true)
	}

	/* Create the disk image to libguestfs. */
	optargsDiskCreate := guestfs.OptargsDisk_create{
		Backingfile_is_set:   false,
		Backingformat_is_set: false,
		Preallocation_is_set: false,
		Compat_is_set:        false,
		Clustersize_is_set:   false,
	}

	/* Create a raw-format sparse disk image, d.DiskSize MB in size. */
	if err := g.Disk_create(output, "raw", int64(d.DiskSize*1024*1024), &optargsDiskCreate); err != nil {
		panic(err)
	}

	/* Attach the disk image to libguestfs. */
	optargsAdd_drive := guestfs.OptargsAdd_drive{
		Format_is_set:   true,
		Format:          "raw",
		Readonly_is_set: true,
		Readonly:        false,
	}

	if err := g.Add_drive(output, &optargsAdd_drive); err != nil {
		panic(err)
	}

	/* Run the libguestfs back-end. */
	if err := g.Launch(); err != nil {
		panic(err)
	}

	/* Get the list of devices.  Because we only added one drive
	 * above, we expect that this list should contain a single
	 * element.
	 */
	devices, err := g.List_devices()
	if err != nil {
		panic(err)
	}
	if len(devices) != 1 {
		panic("expected a single device from list-devices")
	}

	/* Partition the disk as one single MBR partition. */
	err = g.Part_disk(devices[0], "mbr")
	if err != nil {
		panic(err)
	}

	/* Get the list of partitions.  We expect a single element, which
	 * is the partition we have just created.
	 */
	partitions, err := g.List_partitions()
	if err != nil {
		panic(err)
	}
	if len(partitions) != 1 {
		panic("expected a single partition from list-partitions")
	}

	/* Create a filesystem on the partition. */
	err = g.Mkfs("ext4", partitions[0], nil)
	if err != nil {
		panic(err)
	}

	/* Now mount the filesystem so that we can add files. */
	err = g.Mount(partitions[0], "/")
	if err != nil {
		panic(err)
	}

	/* Mkdir -p place of userdata.tar */
	err = g.Mkdir_p("/var/lib/boot2docker")
	if err != nil {
		panic(err)
	}

	/* Uploads the local userdata.tar file into /var/lib/boot2docker */
	err = g.Upload(d.ResolveStorePath("userdata.tar"), "/var/lib/boot2docker/userdata.tar")
	if err != nil {
		panic(err)
	}

	/* Because we wrote to the disk and we want to detect write
	 * errors, call g:shutdown.  You don't need to do this:
	 * g.Close will do it implicitly.
	 */
	if err = g.Shutdown(); err != nil {
		panic(fmt.Sprintf("write to disk failed: %s", err))
	}

	return nil
}

// Make a boot2docker userdata.tar key bundle
func (d *Driver) generateKeyBundle() error { // TODO
	magicString := "boot2docker, this is xhyve speaking"

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// magicString first so the automount script knows to format the disk
	file := &tar.Header{Name: magicString, Size: int64(len(magicString))}
	if err := tw.WriteHeader(file); err != nil {
		return err
	}
	if _, err := tw.Write([]byte(magicString)); err != nil {
		return err
	}
	// .ssh/key.pub => authorized_keys
	file = &tar.Header{Name: ".ssh", Typeflag: tar.TypeDir, Mode: 0700}
	if err := tw.WriteHeader(file); err != nil {
		return err
	}
	pubKey, err := ioutil.ReadFile(d.publicSSHKeyPath())
	if err != nil {
		return err
	}
	file = &tar.Header{Name: ".ssh/authorized_keys", Size: int64(len(pubKey)), Mode: 0644}
	if err := tw.WriteHeader(file); err != nil {
		return err
	}
	if _, err := tw.Write([]byte(pubKey)); err != nil {
		return err
	}
	file = &tar.Header{Name: ".ssh/authorized_keys2", Size: int64(len(pubKey)), Mode: 0644}
	if err := tw.WriteHeader(file); err != nil {
		return err
	}
	if _, err := tw.Write([]byte(pubKey)); err != nil {
		return err
	}
	if err := tw.Close(); err != nil {
		return err
	}
	raw := buf.Bytes()

	if err := ioutil.WriteFile(d.userdataPath(), raw, 0644); err != nil {
		return err
	}

	return nil
}
