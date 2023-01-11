package checkcertfile

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestCertFile(t *testing.T) {
	goroot, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		t.Fatalf("Faild to get GOROOT: %s", err)
	}
	err = exec.Command(
		"go",
		"run",
		strings.TrimSuffix(string(goroot), "\n")+"/src/crypto/tls/generate_cert.go",
		"-host",
		"localhost",
		"-duration",
		"720h0m0s", // 30*24*time.Hour
	).Run()
	if err != nil {
		t.Fatalf("Faild to generate cert.pem: %s", err)
	}

	ckr := Run([]string{"-f", "./cert.pem", "-w", "35", "-c", "25"})
	assert.Equal(t, checkers.WARNING, ckr.Status, "should be WARNING")

	os.Remove("./cert.pem")
	os.Remove("./key.pem")
}
