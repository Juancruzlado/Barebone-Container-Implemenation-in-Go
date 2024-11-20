package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <command> [args...]", os.Args[0])
	}
	switch os.Args[1] {
	case "run":
		run(os.Args[2:])
	case "child":
		child(os.Args[2:])
	default:
		log.Fatalf("Unknown command: %s", os.Args[1])
	}
}

func run(args []string) {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, args...)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr, cmd.SysProcAttr = os.Stdin, os.Stdout, os.Stderr, &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if err := cmd.Start(); err != nil || setupCgroups(cmd.Process.Pid) != nil {
		log.Fatalf("Error: %v", err)
	}
	cmd.Wait()
}

func child(args []string) {
	if err := syscall.Sethostname([]byte("container-hostname")); err != nil ||
		syscall.Chroot("/path/to/new/root") != nil || os.Chdir("/") != nil ||
		syscall.Mount("proc", "/proc", "proc", 0, "") != nil {
		log.Fatalf("Error setting up child environment")
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil || syscall.Unmount("/proc", 0) != nil {
		log.Fatalf("Error running command in child")
	}
}

func setupCgroups(pid int) error {
	for _, group := range []struct {
		path, param, value string
	}{{"/sys/fs/cgroup/cpu/mycontainer", "cpu.cfs_quota_us", "50000"},
		{"/sys/fs/cgroup/memory/mycontainer", "memory.limit_in_bytes", "100000000"}} {
		if err := os.MkdirAll(group.path, 0755); err != nil ||
			os.WriteFile(filepath.Join(group.path, group.param), []byte(group.value), 0644) != nil ||
			os.WriteFile(filepath.Join(group.path, "tasks"), []byte(strconv.Itoa(pid)), 0644) != nil {
			return err
		}
	}
	return nil
}
