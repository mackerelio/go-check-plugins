package checkcertfile

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestCertFile(t *testing.T) {
	tmpDir := t.TempDir()
	cmd := exec.Command(
		"go",
		"run",
		filepath.FromSlash(filepath.Join(runtime.GOROOT(), "/src/crypto/tls/generate_cert.go")),
		"-host",
		"localhost",
		"-duration",
		"720h0m0s", // 30*24*time.Hour
	)
	cmd.Dir = tmpDir
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to generate cert.pem: %s", err)
	}

	ckr := Run([]string{"-f", filepath.Join(tmpDir, "cert.pem"), "-w", "35", "-c", "25"})
	assert.Equal(t, checkers.WARNING, ckr.Status, "should be WARNING")
}
