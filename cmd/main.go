//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("bad command")
	}
}

func run() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func child() {
	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())
	u, _ := user.Current()
	fmt.Println(u.Username)
	cgroups()
	syscall.Sethostname([]byte("container"))
	syscall.Chroot("/ubuntu-fs")
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")
	err := syscall.Exec(os.Args[2], os.Args[2:], os.Environ())
	if err != nil {
		panic("err while execve'ing child process ")
	}
	syscall.Unmount("/proc", 0)
}

func cgroups() {
	cgroupsRoot := "/sys/fs/cgroup/"
	cgroupPath := filepath.Join(cgroupsRoot, "adi")
	err := os.Mkdir(cgroupPath, 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	pid := strconv.Itoa(os.Getpid())
	must(os.WriteFile(filepath.Join(cgroupPath, "pids.max"), []byte("20"), 0700))
	must(os.WriteFile(filepath.Join(cgroupPath, "memory.max"), []byte("52428800"), 0700))
	must(os.WriteFile(filepath.Join(cgroupPath, "memory.swap.max"), []byte("00"), 0700))
	must(os.WriteFile(filepath.Join(cgroupPath, "cgroup.procs"), []byte(pid), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
