package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	prog    string = "ducky"
	reEntry string = "ducky-fork"
	kerUser        = 3.08
	kerUid         = 3.19
)

func check_kernel_ver() (float64, error) {
	// need to get kernel version first
	kstr, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return 0, fmt.Errorf("check_kernel_ver() uname: %s", err)
	}

	idx := strings.LastIndex(string(kstr), ".")
	f, err := strconv.ParseFloat(string(kstr[:idx]), 64)
	if err != nil {
		return 0, fmt.Errorf("check_kernel_ver() parsing float %s. %s",
			kstr[:idx], err)
	}
	if f == 3.2 {
		f = 3.02
	} else if f == 3.3 {
		f = 3.03
	} else if f == 3.4 {
		f = 3.04
	} else if f == 3.5 {
		f = 3.05
	} else if f == 3.6 {
		f = 3.06
	} else if f == 3.7 {
		f = 3.07
	} else if f == 3.8 {
		f = 3.08
	} else if f == 3.9 {
		f = 3.09
	}

	return f, nil
}

func setContainer(cmdName string, args []string) error {
	// need to get kernel version first
	kVer, err := check_kernel_ver()
	if err != nil {
		return fmt.Errorf("setContainer() %s", err)
	}
	cmd := &exec.Cmd{
		Path:   cmdName,
		Args:   append([]string{reEntry}, args...),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	// set the default flags
	var flags uintptr = syscall.CLONE_NEWPID |
		syscall.CLONE_NEWNET |
		syscall.CLONE_NEWUTS |
		syscall.CLONE_NEWNS |
		// using 'ipcmk -Q' to create ipc on host
		// 'ipcs' to verify
		syscall.CLONE_NEWIPC

	if kVer == kerUid {
		// uid mapping is broken in 3.19 kernel
		fmt.Println("kernel 3.19 has a bug in Uid mapping, may not work")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: flags | syscall.CLONE_NEWUSER,
		}
	} else if kVer > kerUser {
		// USER namespace is added after 3.8 kernel
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: flags | syscall.CLONE_NEWUSER,

			UidMappings: []syscall.SysProcIDMap{
				{
					ContainerID: 0,
					HostID:      os.Getuid(),
					Size:        1,
				},
			},
			GidMappings: []syscall.SysProcIDMap{
				{
					ContainerID: 0,
					HostID:      os.Getgid(),
					Size:        1,
				},
			},
		}
	} else {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: flags,
		}
	}
	if err := cmd.Run(); err != nil {
		// fmt.Println(stderr.String())
		// this will return exit status 12 on DDOS ???
		return fmt.Errorf("setContainer() %s: %s ", cmdName, err)
		// return err
	}
	cmd.Wait()
	return nil
}

type Mount struct {
	Source string
	Target string
	Fs     string
	Flags  int
	Data   string
}

type Cfg struct {
	Path     string
	Args     []string
	Hostname string
	Mounts   []Mount
	Rootfs   string
	IP       string
}

var defaultMountFlags = syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

var defaultCfg = Cfg{
	Hostname: "unc",
	Mounts: []Mount{
		{
			Source: "proc",
			Target: "/proc",
			Fs:     "proc",
			Flags:  defaultMountFlags,
		},
		{
			Source: "tmpfs",
			Target: "/dev",
			Fs:     "tmpfs",
			Flags:  syscall.MS_NOSUID | syscall.MS_STRICTATIME,
			Data:   "mode=755",
		},
	},
	Rootfs: "/home/moroz/project/busybox",
}

func mount(cfg Cfg) error {
	for _, m := range cfg.Mounts {
		target := filepath.Join(cfg.Rootfs, m.Target)
		// fmt.Printf("Mount %s to %s\n", m.Source, target)
		if err := syscall.Mount(m.Source, target, m.Fs, uintptr(m.Flags), m.Data); err != nil {
			return fmt.Errorf("failed to mount %s to %s: %v", m.Source, target, err)
		}
	}
	return nil
}

func pivotRoot(root string) error {
	// we need this to satisfy restriction:
	// "new_root and put_old must not be on the same filesystem as the current root"
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount() %v", err)
	}
	// create rootfs/.pivot_root as path for old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return fmt.Errorf("Mkdir() %v", err)
	}
	//fmt.Printf("Pivot root dir: %s\n", pivotDir)
	//fmt.Printf("Pivot root to %s\n", root)
	// pivot_root to rootfs, now old_root is mounted in rootfs/.pivot_root
	// mounts from it still can be seen in `mount`
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("PivotRoot() %v", err)
	}
	// change working directory to /
	// it is recommendation from man-page
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("Chdir(/) %v", err)
	}
	// path to pivot root now changed, update
	pivotDir = filepath.Join("/", ".pivot_root")
	// umount rootfs/.pivot_root(which is now /.pivot_root) with all submounts
	// now we have only mounts that we mounted ourself in `mount`
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("Unmount() %v\n", err)
	}
	// remove temporary directory
	return os.Remove(pivotDir)
}

func fork() error {
	// setup flag
	cmd, err := exec.LookPath(os.Args[1])
	if err != nil {
		return fmt.Errorf("LookPath() %v", err)
	}

	// set the current working dir as rootfs
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Getwd() %v", err)
	}

	defaultCfg.Rootfs = wd
	// store given IP

	// mount
	if err := mount(defaultCfg); err != nil {
		return fmt.Errorf("mount() %v", err)
	}

	// pivot root
	if err := pivotRoot(defaultCfg.Rootfs); err != nil {
		return fmt.Errorf("pivotRoot() %v", err)
	}

	// sethostname()
	if err := syscall.Sethostname([]byte("james.datadomain.com")); err != nil {
		return fmt.Errorf("Sethostname() %v", err)
	}
	// waitfor iface
	// setup net iface
	// create bridge
	//    ip link add name mybridge type bridge
	// create a paire of virtual eth connections between p & c
	//    ip link add name veth0 type veth peer name veth1 netns <pid>
	// https://www.toptal.com/linux/separation-anxiety-isolating-your-system-with-linux-namespaces
	// http://containerops.org/2014/07/30/tenus-golang-powered-linux-networking/
	// http://www.haifux.org/lectures/299/netLec7.pdf

	// exec
	return syscall.Exec(cmd, os.Args[1:], os.Environ())
}

func main() {
	// fmt.Printf("%s: PID %d\n", os.Args[0], syscall.Getpid())

	// this is the 2nd entry.
	if os.Args[0] == reEntry {
		if err := fork(); err != nil {
			fmt.Printf("Error: gons-fork %v\n", err)
		}
		os.Exit(0)
	}

	// this is the 1st entry
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println(reEntry, "<bin> <args> ... ")
		os.Exit(1)
	}

	// setup the proc env and calling clone/exec to enter the 2nd entry
	if err := setContainer(os.Args[0], os.Args[1:]); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
