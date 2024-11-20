package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestRunCommand(t *testing.T) {
	// Use a simple command for testing (e.g., `echo`)
	cmd := exec.Command("/proc/self/exe", "run", "echo", "test")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Run command failed: %v", err)
	}
}

func TestChildCommand(t *testing.T) {
	// Test the child process execution by running a harmless command
	cmd := exec.Command("/proc/self/exe", "child", "echo", "child process test")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Child command failed: %v", err)
	}

	if !strings.Contains(string(output), "child process test") {
		t.Fatalf("Expected output not found: %s", string(output))
	}
}

func TestSetupCgroups(t *testing.T) {
	// Create a temporary directory to mock cgroup paths
	tempDir := t.TempDir()
	cgroupPath := filepath.Join(tempDir, "cgroup")
	os.Setenv("CGROUP_PATH", cgroupPath) // Mock the cgroup root

	defer os.Unsetenv("CGROUP_PATH")

	// Use a mock PID for testing
	mockPid := 12345

	if err := setupCgroups(mockPid); err != nil {
		t.Fatalf("Cgroup setup failed: %v", err)
	}

	// Check that cgroup files were created with the correct values
	cpuQuota, err := os.ReadFile(filepath.Join(cgroupPath, "cpu", "mycontainer", "cpu.cfs_quota_us"))
	if err != nil || string(cpuQuota) != "50000" {
		t.Fatalf("CPU quota setup failed")
	}

	memLimit, err := os.ReadFile(filepath.Join(cgroupPath, "memory", "mycontainer", "memory.limit_in_bytes"))
	if err != nil || string(memLimit) != "100000000" {
		t.Fatalf("Memory limit setup failed")
	}

	tasksFile, err := os.ReadFile(filepath.Join(cgroupPath, "cpu", "mycontainer", "tasks"))
	if err != nil || string(tasksFile) != strconv.Itoa(mockPid) {
		t.Fatalf("Tasks assignment failed")
	}
}
