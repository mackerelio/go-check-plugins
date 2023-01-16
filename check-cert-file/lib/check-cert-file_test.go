package checkcertfile

import (
	// "os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestCertFile(t *testing.T) {
	goroot, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		t.Fatalf("Failed to get GOROOT: %s", err)
	}

	tmpDir := t.TempDir()
	cmd := exec.Command(
		"go",
		"run",
		filepath.FromSlash(filepath.Join(strings.TrimSuffix(string(goroot), "\n"), "/src/crypto/tls/generate_cert.go")),
		"-host",
		"localhost",
		"-duration",
		"720h0m0s", // 30*24*time.Hour
	)
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to generate cert.pem: %s", err)
	}

	ckr := Run([]string{"-f", filepath.Join(tmpDir, "cert.pem"), "-w", "35", "-c", "25"})
	assert.Equal(t, checkers.WARNING, ckr.Status, "should be WARNING")
}
