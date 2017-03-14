package checkntservice

import (
	"os/exec"
	"runtime"
	"testing"
)

func stopSpoolerService() error {
	_, err := exec.Command("net", "stop", "Spooler").CombinedOutput()
	return err
}

func startSpoolerService() error {
	_, err := exec.Command("net", "start", "Spooler").CombinedOutput()
	return err
}

func TestNtService(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip(runtime.GOOS + " doesn't implement Windows NT service")
	}

	ss, err := getServiceState()
	if err != nil {
		t.Errorf("failed to get service status: %v", err)
	}
	for _, s := range ss {
		if s.Name == "Spooler" {
			if s.State != "Running" {
				t.Error("Spooler service should be started in default")
			}
		}
	}

	err = stopSpoolerService()
	if err != nil {
		t.Skipf("failed to stop Spooler service. But ignore this: %v", err)
	}
	defer startSpoolerService()

	ss, err = getServiceState()
	if err != nil {
		t.Errorf("failed to get service status: %v", err)
	}
	for _, s := range ss {
		if s.Name == "Spooler" {
			if s.State == "Running" {
				t.Error("Spooler service should be stopped now")
			}
		}
	}
}
