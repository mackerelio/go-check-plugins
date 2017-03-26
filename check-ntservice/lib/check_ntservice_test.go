package checkntservice

import (
	"os/exec"
	"runtime"
	"syscall"
	"testing"
)

func stopFaxService() error {
	_, err := exec.Command("net", "stop", "Fax").CombinedOutput()
	return err
}

func startFaxService() error {
	_, err := exec.Command("net", "start", "Fax").CombinedOutput()
	return err
}

func TestNtService(t *testing.T) {
	ss, err := getServiceState()
	if runtime.GOOS != "windows" {
		if err == nil || err != syscall.ENOSYS {
			t.Fatal(runtime.GOOS + " should fail because it's not Windows")
		}
		t.Skip(runtime.GOOS + " doesn't implement Windows NT service")
	}
	if err != nil {
		t.Errorf("failed to get service status: %v", err)
	}
	for _, s := range ss {
		if s.Name == "Fax" {
			if s.State != "Running" {
				t.Errorf("Fax service should be started in default: %v", s.State)
			}
		}
	}

	err = stopFaxService()
	if err != nil {
		t.Skipf("failed to stop Fax service. But ignore this: %v", err)
	}
	defer startFaxService()

	ss, err = getServiceState()
	if err != nil {
		t.Errorf("failed to get service status: %v", err)
	}
	for _, s := range ss {
		if s.Name == "Fax" {
			if s.State == "Running" {
				t.Error("Fax service should be stopped now")
			}
		}
	}
}
