// Copyright 2015 The docker-machine-driver-xhyve Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xhyve

import (
	"archive/tar"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
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
	ps "github.com/mitchellh/go-ps"
	hyperkit "github.com/moby/hyperkit/go"
	"github.com/zchee/docker-machine-driver-xhyve/b2d"
	"github.com/zchee/docker-machine-driver-xhyve/vmnet"
	qcow2 "github.com/zchee/go-qcow2"
	govmnet "github.com/zchee/go-vmnet"
)

const (
	isoFilename           = "boot2docker.iso"
	isoMountPath          = "b2d-image"
	defaultBootKernel     = ""
	defaultBootInitrd     = ""
	defaultBoot2DockerURL = ""
	defaultBootCmd        = ""
	defaultCPU            = 1
	defaultCaCertPath     = ""
	defaultDiskSize       = 20000
	defaultMacAddr        = ""
	defaultMemory         = 1024
	defaultISOFilename    = "boot2docker.iso"
	defaultPrivateKeyPath = ""
	defaultUUID           = ""
	defaultNFSShare       = false
	rootVolumeName        = "root-volume"
	defaultDiskNumber     = -1
	defaultVirtio9p       = false
	defaultQcow2          = false
	defaultRawDisk        = false
)

type Driver struct {
	*drivers.BaseDriver
	*b2d.B2dUtils

	Boot2DockerURL string
	CaCertPath     string
	PrivateKeyPath string

	CPU            int
	Memory         int
	DiskSize       int64
	DiskNumber     int
	MacAddr        string
	UUID           string
	Qcow2          bool
	RawDisk        bool
	Virtio9p       bool
	Virtio9pFolder string
	NFSShare       bool

	BootCmd    string
	BootKernel string
	BootInitrd string
	Initrd     string
	Vmlinuz    string
}

var (
	ErrMachineExist    = errors.New("machine already exists")
	ErrMachineNotExist = errors.New("machine does not exist")
	diskRegexp         = regexp.MustCompile("^/dev/disk([0-9]+)")
	kernelRegexp       = regexp.MustCompile(`(vmlinu[xz]|bzImage)[\d]*`)
	kernelOptionRegexp = regexp.MustCompile(`(?:\t|\s{2})append\s+([[:print:]]+)`)
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
		BootKernel:     defaultBootKernel,
		BootInitrd:     defaultBootInitrd,
		CPU:            defaultCPU,
		CaCertPath:     defaultCaCertPath,
		DiskSize:       defaultDiskSize,
		MacAddr:        defaultMacAddr,
		Memory:         defaultMemory,
		PrivateKeyPath: defaultPrivateKeyPath,
		UUID:           defaultUUID,
		NFSShare:       defaultNFSShare,
		DiskNumber:     defaultDiskNumber,
		Virtio9p:       defaultVirtio9p,
		Qcow2:          defaultQcow2,
		RawDisk:        defaultRawDisk,
	}
}

// RegisterCreateFlags registers the flags this driver adds to
// "docker hosts create"
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "XHYVE_BOOT_CMD",
			Name:   "xhyve-boot-cmd",
			Usage:  "Command of booting kexec protocol",
			Value:  defaultBootCmd,
		},
		mcnflag.StringFlag{
			EnvVar: "XHYVE_BOOT_KERNEL",
			Name:   "xhyve-boot-kernel",
			Usage:  "Absolute path to kernel file (like /boot/vmlinuz64)",
			Value:  defaultBootKernel,
		},
		mcnflag.StringFlag{
			EnvVar: "XHYVE_BOOT_INITRD",
			Name:   "xhyve-boot-initrd",
			Usage:  "Absolute path to ramdisk file (like /boot/initrd.img)",
			Value:  defaultBootInitrd,
		},
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
			EnvVar: "XHYVE_DISK_SIZE",
			Name:   "xhyve-disk-size",
			Usage:  "Size of disk for host in MB",
			Value:  defaultDiskSize,
		},
		mcnflag.BoolFlag{
			EnvVar: "XHYVE_EXPERIMENTAL_NFS_SHARE",
			Name:   "xhyve-experimental-nfs-share",
			Usage:  "Setup NFS shared folder (requires root)",
		},
		mcnflag.IntFlag{
			EnvVar: "XHYVE_MEMORY_SIZE",
			Name:   "xhyve-memory-size",
			Usage:  "Size of memory for host in MB",
			Value:  defaultMemory,
		},
		mcnflag.BoolFlag{
			EnvVar: "XHYVE_QCOW2",
			Name:   "xhyve-qcow2",
			Usage:  "Use qcow2 disk format",
		},
		mcnflag.BoolFlag{
			EnvVar: "XHYVE_RAW_DISK",
			Name:   "xhyve-rawdisk",
			Usage:  "Use a raw disk for attached volumes",
		},
		mcnflag.StringFlag{
			EnvVar: "XHYVE_UUID",
			Name:   "xhyve-uuid",
			Usage:  "The UUID for the machine",
			Value:  defaultUUID,
		},
		mcnflag.BoolFlag{
			EnvVar: "XHYVE_VIRTIO_9P",
			Name:   "xhyve-virtio-9p",
			Usage:  "Setup virtio-9p folder share",
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
	d.BootCmd = flags.String("xhyve-boot-cmd")
	d.BootKernel = flags.String("xhyve-boot-kernel")
	d.BootInitrd = flags.String("xhyve-boot-initrd")
	d.CPU = flags.Int("xhyve-cpu-count")
	if d.CPU < 1 {
		d.CPU = int(runtime.NumCPU())
	}
	d.DiskSize = int64(flags.Int("xhyve-disk-size"))
	d.Memory = flags.Int("xhyve-memory-size")
	d.NFSShare = flags.Bool("xhyve-experimental-nfs-share")
	d.Qcow2 = flags.Bool("xhyve-qcow2")
	d.RawDisk = flags.Bool("xhyve-rawdisk")
	d.SSHPort = 22
	d.SSHUser = "docker"
	d.SwarmDiscovery = flags.String("swarm-discovery")
	d.SwarmHost = flags.String("swarm-host")
	d.SwarmMaster = flags.Bool("swarm-master")
	d.UUID = flags.String("xhyve-uuid")
	d.Virtio9p = flags.Bool("xhyve-virtio-9p")
	d.Virtio9pFolder = "/Users"

	return nil
}

// PreCommandCheck Check required of docker-machine-driver-xhyve before any func
// func: GetURL, PreCreateCheck, Start, Stop, Restart
func (d *Driver) PreCommandCheck() error {
	bin, err := os.Stat(os.Args[0])
	if err != nil {
		return err
	}

	// Check of own binary owner and uid
	if int(bin.Sys().(*syscall.Stat_t).Uid) != 0 {
		return fmt.Errorf("%s binary needs root owner and uid. See https://github.com/zchee/docker-machine-driver-xhyve#install", bin.Name())
	}

	// Check of execute user
	user := syscall.Getuid()
	if user == 0 {
		return fmt.Errorf("%s needs to be executed with the privileges of the user. please remove sudo on execute command", bin.Name())
	}

	return nil
}

func (d *Driver) GetURL() (string, error) {
	if err := d.PreCommandCheck(); err != nil {
		return "", err
	}

	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	if ip == "" {
		return "", nil
	}

	// Wait for SSH over NAT to be available before returning to user
	for {
		err := drivers.WaitForSSH(d)
		if err != nil {
			time.Sleep(1 * time.Second)
		} else {
			break
		}
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

	if d.IPAddress != "" {
		return d.IPAddress, nil
	}

	return d.getIPfromDHCPLease()
}

func (d *Driver) GetState() (state.State, error) {
	pid, err := d.GetPid()
	if err != nil {
		// TODO: If err instead of nil, will be occurred error when first GetState() of Start()
		return state.Error, nil
	}

	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return state.Error, err
	}

	if err := proc.Signal(syscall.Signal(0)); err != nil {
		return state.Stopped, nil
	}

	psproc, err := ps.FindProcess(int(pid))
	if err != nil {
		return state.Error, err
	}
	// process name is truncated to 'docker-machine-d'
	exe := psproc.Executable()
	if !(strings.Contains(exe, "docker-machine") || strings.Contains(exe, "hyperkit")) {
		return state.Error, fmt.Errorf("Unable to find 'xhyve' process by PID: %d", pid)
	}

	return state.Running, nil
}

func (d *Driver) waitForIP() error {
	var ip string
	var err error

	log.Infof("Waiting for VM to come online...")
	for i := 1; i <= 60; i++ {
		ip, err = d.getIPfromDHCPLease()
		if err != nil {
			log.Debugf("Not there yet %d/%d, error: %s", i, 60, err)
			time.Sleep(2 * time.Second)
			continue
		}

		if ip != "" {
			log.Debugf("Got an ip: %s", ip)
			d.IPAddress = ip

			break
		}
	}

	if ip == "" {
		return fmt.Errorf("Machine didn't return an IP after 120 seconds, aborting")
	}

	// Wait for SSH over NAT to be available before returning to user
	if err := drivers.WaitForSSH(d); err != nil {
		return err
	}

	return nil
}

// PreCreateCheck Prints driver version, and Check VirtualBox version
func (d *Driver) PreCreateCheck() error {
	// Check required of docker-machine-driver-xhyve
	if err := d.PreCommandCheck(); err != nil {
		return err
	}

	//TODO: libmachine PLEASE output driver version API!
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
	if err := d.CopyIsoToMachineDir(d.Boot2DockerURL, d.MachineName); err != nil {
		return err
	}

	log.Infof("Creating VM...")
	if err := os.MkdirAll(d.ResolveStorePath("."), 0755); err != nil {
		return err
	}

	if err := d.extractKernelImages(); err != nil {
		return err
	}

	log.Infof("Generating %dMB disk image...", d.DiskSize)

	if d.Qcow2 {
		if err := d.generateQcow2Image(d.DiskSize); err != nil {
			return err
		}
	} else if d.RawDisk {
		if err := d.generateRawDiskImage(d.DiskSize); err != nil {
			return err
		}
	} else {
		if err := d.generateSparseBundleDiskImage(d.DiskSize); err != nil {
			return err
		}
	}

	// Fix file permission root to current user for vmnet.framework
	log.Infof("Fix file permission...")
	os.Chown(d.ResolveStorePath("."), syscall.Getuid(), syscall.Getegid())
	files, _ := ioutil.ReadDir(d.ResolveStorePath("."))
	for _, f := range files {
		log.Debugf(d.ResolveStorePath(f.Name()))
		os.Chown(d.ResolveStorePath(f.Name()), syscall.Getuid(), syscall.Getegid())
	}

	if d.UUID == "" {
		log.Infof("Generate UUID...")
		d.UUID = uuidgen()
		log.Debugf("Generated UUID: %s", d.UUID)
	} else {
		log.Infof("Using Supplied UUID: %s", d.UUID)
	}

	log.Infof("Convert UUID to MAC address...")
	rawUUID, err := d.getMACAdress()
	if err != nil {
		return fmt.Errorf("Could not convert the UUID to MAC address: %s", err.Error())
	}
	d.MacAddr = trimMacAddress(rawUUID)
	log.Debugf("Converted MAC address: %s", d.MacAddr)

	log.Infof("Starting %s...", d.MachineName)
	if err := d.Start(); err != nil {
		return err
	}

	d.setupMounts()

	return nil
}

func (d *Driver) Start() error {
	if err := d.PreCommandCheck(); err != nil {
		return err
	}

	pid := d.ResolveStorePath(d.MachineName + ".pid")
	if _, err := os.Stat(pid); err == nil {
		os.Remove(pid)
	}

	d.attachDiskImage()

	h, err := hyperkit.New("", "auto", d.ResolveStorePath(""))
	if err != nil {
		return err
	}
	h.Console = hyperkit.ConsoleFile

	h.Kernel = d.ResolveStorePath(d.Vmlinuz)
	h.Initrd = d.ResolveStorePath(d.Initrd)

	if d.Qcow2 {
		h.Disks = []hyperkit.DiskConfig{
			{
				Path:   fmt.Sprintf("file://%s", filepath.Join(d.ResolveStorePath("."), d.MachineName+".qcow2")),
				Format: "qcow",
				Driver: "virtio-blk",
			},
		}
	} else if d.RawDisk {
		h.Disks = []hyperkit.DiskConfig{
			{
				Path:   filepath.Join(d.ResolveStorePath("."), d.MachineName+".rawdisk"),
				Driver: "virtio-blk",
			},
		}
	} else {
		h.Disks = []hyperkit.DiskConfig{
			{
				Path:   fmt.Sprintf("/dev/rdisk%d", d.DiskNumber),
				Driver: "ahci-hd",
			},
		}
	}

	h.UUID = d.UUID
	h.CPUs = d.CPU
	h.Memory = d.Memory
	h.VMNet = true

	if err := h.Start(d.BootCmd); err != nil {
		return err
	}

	if err := d.waitForIP(); err != nil {
		return err
	}

	if err := d.setupMounts(); err != nil {
		log.Warnf("Error setting up mounts: %v", err)
	}

	return nil
}

func (d *Driver) Stop() error {
	if err := d.PreCommandCheck(); err != nil {
		return err
	}

	log.Infof("Stopping %s ...", d.MachineName)
	if err := d.SendSignal(syscall.SIGTERM); err != nil {
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
	d.detachDiskImage()

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
	}

	if err := d.removeDiskImage(); err != nil {
		return err
	}

	if d.NFSShare {
		log.Infof("Remove NFS share folder must be root. Please insert root password.")
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
	if err := d.PreCommandCheck(); err != nil {
		return err
	}

	s, err := d.GetState()
	if err != nil {
		return err
	}
	if s == state.Running {
		if err := d.Stop(); err != nil {
			return err
		}
	}

	if err := d.Start(); err != nil {
		return err
	}

	return d.waitForIP()
}

func (d *Driver) Kill() error {
	log.Infof("Killing %s ...", d.MachineName)
	if err := d.SendSignal(syscall.SIGKILL); err != nil {
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

func readLine(path string) (string, error) {
	inFile, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		if kernelOptionRegexp.Match(scanner.Bytes()) {
			m := kernelOptionRegexp.FindSubmatch(scanner.Bytes())
			return string(m[1]), nil
		}
	}
	return "", fmt.Errorf("couldn't find kernel option from %s image", path)
}

func (d *Driver) extractKernelOptions() error {
	volumeRootDir := d.ResolveStorePath(isoMountPath)
	if d.BootCmd == "" {
		err := filepath.Walk(volumeRootDir, func(path string, f os.FileInfo, err error) error {
			if strings.Contains(path, "isolinux.cfg") {
				d.BootCmd, err = readLine(path)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}

		if d.BootCmd == "" {
			return errors.New("Not able to parse isolinux.cfg, Please use --xhyve-boot-cmd option")
		}
	}

	log.Debugf("Extracted Options %q", d.BootCmd)
	return nil
}

func (d *Driver) extractKernelImages() error {
	log.Debugf("Mounting %s", isoFilename)

	volumeRootDir := d.ResolveStorePath(isoMountPath)
	err := hdiutil("attach", d.ResolveStorePath(isoFilename), "-mountpoint", volumeRootDir)
	if err != nil {
		return err
	}

	log.Debugf("Extracting Kernel Options...")
	if err := d.extractKernelOptions(); err != nil {
		return err
	}

	defer func() error {
		log.Debugf("Unmounting %s", isoFilename)
		return hdiutil("detach", volumeRootDir)
	}()

	if d.BootKernel == "" && d.BootInitrd == "" {
		err = filepath.Walk(volumeRootDir, func(path string, f os.FileInfo, err error) error {
			if kernelRegexp.MatchString(path) {
				d.BootKernel = path
				_, d.Vmlinuz = filepath.Split(path)
			}
			if strings.Contains(path, "initrd") {
				d.BootInitrd = path
				_, d.Initrd = filepath.Split(path)
			}
			return nil
		})
	}

	if err != nil {
		if err != nil || d.BootKernel == "" || d.BootInitrd == "" {
			err = fmt.Errorf("==== Can't extract Kernel and Ramdisk file ====")
			return err
		}
	}

	dest := d.ResolveStorePath(d.Vmlinuz)
	log.Debugf("Extracting %s into %s", d.BootKernel, dest)
	if err := mcnutils.CopyFile(d.BootKernel, dest); err != nil {
		return err
	}

	dest = d.ResolveStorePath(d.Initrd)
	log.Debugf("Extracting %s into %s", d.BootInitrd, dest)
	if err := mcnutils.CopyFile(d.BootInitrd, dest); err != nil {
		return err
	}

	return nil
}

func (d *Driver) generateRawDiskImage(size int64) error {
	diskPath := filepath.Join(d.ResolveStorePath("."), d.MachineName+".rawdisk")

	f, err := os.OpenFile(diskPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	f.Close()

	if err := os.Truncate(diskPath, d.DiskSize*1048576); err != nil {
		return err
	}

	tarBuf, err := d.generateKeyBundle()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(diskPath, os.O_WRONLY, 0644)
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

func (d *Driver) generateQcow2Image(size int64) error {
	diskPath := filepath.Join(d.ResolveStorePath("."), d.MachineName+".qcow2")
	opts := &qcow2.Opts{
		Filename:      diskPath,
		Size:          d.DiskSize * 107374,
		Fmt:           qcow2.DriverQCow2,
		ClusterSize:   65536,
		Preallocation: qcow2.PREALLOC_MODE_OFF,
		Encryption:    false,
		LazyRefcounts: true,
	}

	img, err := qcow2.Create(opts)
	if err != nil {
		log.Error(err)
	}

	tarBuf, err := d.generateKeyBundle()
	if err != nil {
		return err
	}

	// TODO(zchee): hardcoded
	zeroFill(tarBuf, 109569)
	tarBuf.Write(diskimageFooter)

	// TODO(zchee): hardcoded
	zeroFill(tarBuf, 16309)
	tarBuf.Write(efipartFooter)

	img.Write(tarBuf.Bytes())

	return nil
}

func (d *Driver) setupMounts() error {
	if d.Virtio9p {
		err := d.setupVirt9pShare()
		if err != nil {
			log.Errorf("virtio-9p setup failed: %s", err.Error())
			return err
		}
	}

	// Setup NFS sharing
	if d.NFSShare {
		log.Infof("NFS share folder must be root. Please insert root password.")
		err := d.setupNFSShare()
		if err != nil {
			log.Errorf("NFS setup failed: %s", err.Error())
			return err
		}
	}
	return nil
}

// zeroFill writes n zero bytes into w.
func zeroFill(w io.Writer, n int64) error {
	const blocksize = 32 << 10
	zeros := make([]byte, blocksize)
	var k int
	var err error
	for n > 0 {
		if n > blocksize {
			k, err = w.Write(zeros)
			if err != nil {
				return err
			}

		} else {
			k, err = w.Write(zeros[:n])
			if err != nil {
				return err
			}

		}
		if err != nil {
			return err
		}
		n -= int64(k)
	}
	return nil
}

func (d *Driver) generateSparseBundleDiskImage(count int64) error {
	diskPath := d.ResolveStorePath(rootVolumeName + ".sparsebundle")

	if err := hdiutil("create", "-megabytes", fmt.Sprintf("%d", count), "-type", "SPARSEBUNDLE", diskPath); err != nil {
		return err
	}

	if err := d.attachDiskImage(); err != nil {
		return err
	}

	tarBuf, err := d.generateKeyBundle()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fmt.Sprintf("/dev/rdisk%d", d.DiskNumber), os.O_WRONLY, 0644)
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

func (d *Driver) attachDiskImage() error {
	diskPath := d.ResolveStorePath(rootVolumeName + ".sparsebundle")
	cmd := exec.Command("hdiutil", "attach", "-nomount", "-noverify", "-noautofsck", diskPath)
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	matches := diskRegexp.FindSubmatch(output)
	if len(matches) != 2 {
		return fmt.Errorf("Failed parsing disk number, hdiutil output: %s", string(output))
	}

	d.DiskNumber, err = strconv.Atoi(string(matches[1]))
	if err != nil {
		return err
	}

	return nil
}

func (d *Driver) detachDiskImage() error {
	if err := hdiutil("detach", fmt.Sprintf("/dev/disk%d", d.DiskNumber)); err != nil {
		return err
	}

	d.DiskNumber = -1

	return nil
}

func (d *Driver) removeDiskImage() error {
	diskPath := d.ResolveStorePath(rootVolumeName + ".sparsebundle")
	return os.RemoveAll(diskPath)
}

// Make a boot2docker userdata.tar key bundle
func (d *Driver) generateKeyBundle() (*bytes.Buffer, error) {
	magicString := "boot2docker, please format-me"

	log.Infof("Creating SSH key...")
	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return nil, err
	}

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

func (d *Driver) setupVirt9pShare() error {
	user, err := user.Current()
	if err != nil {
		return err
	}
	bootScriptName := "/var/lib/boot2docker/bootlocal.sh"
	bootScript := fmt.Sprintf("#/bin/bash\\n"+
		"sudo mkdir -p %s\\n"+
		"sudo mount -t 9p -o version=9p2000 -o trans=virtio -o uname=%s -o dfltuid=$(id -u docker) -o dfltgid=50 -o access=any host %s", d.Virtio9pFolder, user.Username, d.Virtio9pFolder)

	writeScriptCmd := fmt.Sprintf("echo -e \"%s\" | sudo tee %s && sudo chmod +x %s && %s",
		bootScript, bootScriptName, bootScriptName, bootScriptName)

	if _, err := drivers.RunSSHCommandFromDriver(d, writeScriptCmd); err != nil {
		return err
	}

	return nil
}

// Setup NFS share
func (d *Driver) setupNFSShare() error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	nfsConfig := fmt.Sprintf("/Users %s -alldirs -mapall=%s", d.IPAddress, user.Username)

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

func (d *Driver) GetPid() (int, error) {
	p, err := ioutil.ReadFile(d.ResolveStorePath("hyperkit.pid"))
	if err != nil {
		return 0, err
	}

	pid, err := strconv.ParseInt(string(p), 10, 64)
	if err != nil {
		return 0, err
	}

	return int(pid), nil
}

func (d *Driver) SendSignal(sig os.Signal) error {
	pid, err := d.GetPid()
	if err != nil {
		return err
	}

	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return err
	}

	if err := proc.Signal(sig); err != nil {
		return err
	}

	return nil
}

// trimMacAddress trimming "0" of the ten's digit
func trimMacAddress(rawUUID string) string {
	re := regexp.MustCompile(`0([A-Fa-f0-9](:|$))`)
	mac := re.ReplaceAllString(rawUUID, "$1")

	return mac
}

func (d *Driver) getMACAdress() (string, error) {
	m, err := govmnet.GetMACAddressFromUUID(d.UUID)
	if err != nil {
		return "", err
	}

	hw, err := net.ParseMAC(string(m))
	if err != nil {
		return "", err
	}
	return hw.String(), nil
}

func (d *Driver) UpdateISOCache(isoURL string) error {
	b2d := b2d.NewB2dUtils(d.StorePath)
	mcnutils := mcnutils.NewB2dUtils(d.StorePath)

	// recreate the cache dir if it has been manually deleted
	if _, err := os.Stat(b2d.ImgCachePath); os.IsNotExist(err) {
		log.Infof("Image cache directory does not exist, creating it at %s...", b2d.ImgCachePath)
		if err := os.Mkdir(b2d.ImgCachePath, 0700); err != nil {
			return err
		}
	}

	// Check owner of storage cache directory
	cacheStat, _ := os.Stat(b2d.ImgCachePath)

	if int(cacheStat.Sys().(*syscall.Stat_t).Uid) == 0 {
		log.Debugf("Fix %s directory permission...", cacheStat.Name())
		os.Chown(b2d.ImgCachePath, syscall.Getuid(), syscall.Getegid())
	}

	if isoURL != "" {
		// Non-default B2D are not cached
		return nil
	}

	exists := b2d.Exists()
	if !exists {
		log.Info("No default Boot2Docker ISO found locally, downloading the latest release...")
		return mcnutils.DownloadLatestBoot2Docker("")
	}

	latest := b2d.IsLatest()
	if !latest {
		log.Info("Default Boot2Docker ISO is out-of-date, downloading the latest release...")
		return mcnutils.DownloadLatestBoot2Docker("")
	}

	return nil
}

func (d *Driver) CopyIsoToMachineDir(isoURL, machineName string) error {
	b2d := b2d.NewB2dUtils(d.StorePath)
	mcnutils := mcnutils.NewB2dUtils(d.StorePath)

	if err := d.UpdateISOCache(isoURL); err != nil {
		return err
	}

	isoPath := filepath.Join(b2d.ImgCachePath, isoFilename)
	if isoStat, err := os.Stat(isoPath); err == nil {
		if int(isoStat.Sys().(*syscall.Stat_t).Uid) == 0 {
			log.Debugf("Fix %s file permission...", isoStat.Name())
			os.Chown(isoPath, syscall.Getuid(), syscall.Getegid())
		}
	}

	// TODO: This is a bit off-color.
	machineDir := filepath.Join(d.StorePath, "machines", machineName)
	machineIsoPath := filepath.Join(machineDir, isoFilename)

	// By default just copy the existing "cached" iso to the machine's directory...
	defaultISO := filepath.Join(b2d.ImgCachePath, defaultISOFilename)
	if isoURL == "" {
		log.Infof("Copying %s to %s...", defaultISO, machineIsoPath)
		return CopyFile(defaultISO, machineIsoPath)
	}

	// if ISO is specified, check if it matches a github releases url or fallback to a direct download
	downloadURL, err := b2d.GetReleaseURL(isoURL)
	if err != nil {
		return err
	}

	return mcnutils.DownloadISO(machineDir, b2d.Filename(), downloadURL)
}
