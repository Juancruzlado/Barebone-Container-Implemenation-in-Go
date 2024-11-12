package main

import (
	"fmt"
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
	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "run":
		run(args)
	case "child":
		child(args)
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

// run sets up the initial environment and starts a new process in isolated namespaces.
func run(args []string) {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWIPC,
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Error starting command: %v", err)
	}

	// Setup cgroups for the new process
	if err := setupCgroups(cmd.Process.Pid); err != nil {
		log.Fatalf("Error setting up cgroups: %v", err)
	}

	cmd.Wait()
}

// child is the function that runs in the isolated namespace.
func child(args []string) {
	// Set hostname in the new UTS namespace
	if err := syscall.Sethostname([]byte("container-hostname")); err != nil {
		log.Fatalf("Error setting hostname: %v", err)
	}

	// Setup chroot for filesystem isolation
	if err := syscall.Chroot("/path/to/new/root"); err != nil {
		log.Fatalf("Error with chroot: %v", err)
	}
	if err := os.Chdir("/"); err != nil {
		log.Fatalf("Error changing directory: %v", err)
	}

	// Mount proc filesystem
	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		log.Fatalf("Error mounting /proc: %v", err)
	}

	// Execute the desired command
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error running command in child: %v", err)
	}

	// Unmount proc filesystem when done
	if err := syscall.Unmount("/proc", 0); err != nil {
		log.Fatalf("Error unmounting /proc: %v", err)
	}
}

// setupCgroups sets up cgroups to limit resources for the containerized process.
func setupCgroups(pid int) error {
	cgroupPath := "/sys/fs/cgroup/"
	cpu := filepath.Join(cgroupPath, "cpu", "mycontainer")
	mem := filepath.Join(cgroupPath, "memory", "mycontainer")

	// Create cgroup directories
	if err := os.MkdirAll(cpu, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(mem, 0755); err != nil {
		return err
	}

	// Set CPU and memory limits
	if err := os.WriteFile(filepath.Join(cpu, "cpu.cfs_quota_us"), []byte("50000"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(mem, "memory.limit_in_bytes"), []byte("100000000"), 0644); err != nil {
		return err
	}

	// Add the process to the cgroup
	if err := os.WriteFile(filepath.Join(cpu, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(mem, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return err
	}

	return nil
}
