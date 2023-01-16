package checksslcert

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

const (
	gen_cert_path   = "./testdata/gen-cert.sh"
	gen_cert_script = "./gen-cert.sh"
	extfile_path    = "./testdata/extfile.txt"
	ca_crt          = "ca.crt"
	ca_crt_         = "ca.crt"
	client_crt      = "client.crt"
	client_key      = "client.key"
	server_crt      = "server.crt"
	server_key      = "server.key"
)

func prepareCertification(t *testing.T) (string, error) {
	tmpDir := t.TempDir()
	err := exec.Command("cp", gen_cert_path, extfile_path, tmpDir).Run()
	if err != nil {
		return "", fmt.Errorf("Failed to cp: %s", err)
	}

	cm := exec.Command("chmod", "777", gen_cert_script)
	cm.Dir = tmpDir
	err = cm.Run()
	if err != nil {
		return "", fmt.Errorf("Failed to chmod: %s", err)
	}

	cmd := exec.Command("sh", "-c", gen_cert_script)
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Failed to run: %s", err)
	}

	return tmpDir, nil
}

func TestSelfCertification(t *testing.T) {
	tmpDir, err := prepareCertification(t)
	if err != nil {
		t.Errorf("Failed to prepare: %s", err)
	}

	mux := http.NewServeMux()
	ts := httptest.NewUnstartedServer(mux)

	cert, err := tls.LoadX509KeyPair(filepath.Join(tmpDir, server_crt), filepath.Join(tmpDir, server_key))
	if err != nil {
		t.Errorf("Failed to LoadX509KeyPair: %s", err)
	}
	ts.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}

	ts.StartTLS()
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	host, port, _ := net.SplitHostPort(u.Host)

	ckr := Run([]string{
		"-H", host, "-p",
		port, "--ca-file", filepath.Join(tmpDir, ca_crt),
		"-c", "25",
		"-w", "30"})
	assert.Equal(t, checkers.WARNING, ckr.Status, "should be WARNING")
}

func TestClientCertification(t *testing.T) {
	tmpDir, err := prepareCertification(t)
	if err != nil {
		t.Errorf("Failed to prepare: %s", err)
	}

	mux := http.NewServeMux()
	ts := httptest.NewUnstartedServer(mux)

	caCertPEM, err := os.ReadFile(filepath.Join(tmpDir, ca_crt))
	if err != nil {
		t.Errorf("Failed to read file: %s", err)
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caCertPEM)
	if !ok {
		panic("failed to parse root certificate")
	}

	cert, err := tls.LoadX509KeyPair(filepath.Join(tmpDir, server_crt), filepath.Join(tmpDir, server_key))
	if err != nil {
		t.Errorf("Failed to LoadX509KeyPair: %s", err)
	}
	ts.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    roots,
	}

	ts.StartTLS()
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	host, port, _ := net.SplitHostPort(u.Host)

	ckr := Run([]string{
		"-H", host,
		"-p", port,
		"--ca-file", filepath.Join(tmpDir, ca_crt),
		"--cert-file", filepath.Join(tmpDir, client_crt),
		"--key-file", filepath.Join(tmpDir, client_key),
		"-c", "25",
		"-w", "30"})
	assert.Equal(t, checkers.WARNING, ckr.Status, "should be WARNING")
}

func TestNoCheckCertificate(t *testing.T) {
	tmpDir, err := prepareCertification(t)
	if err != nil {
		t.Errorf("Failed to prepare: %s", err)
	}

	mux := http.NewServeMux()
	ts := httptest.NewUnstartedServer(mux)

	cert, err := tls.LoadX509KeyPair(filepath.Join(tmpDir, server_crt), filepath.Join(tmpDir, server_key))
	if err != nil {
		t.Errorf("Failed to LoadX509KeyPair: %s", err)
	}
	ts.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}

	ts.StartTLS()

	u, _ := url.Parse(ts.URL)
	host, port, _ := net.SplitHostPort(u.Host)

	ckr := Run([]string{"-H", host, "-p", port, "-c", "25", "-w", "30"})
	assert.Equal(t, checkers.CRITICAL, ckr.Status, "should be CRITICAL")
	assert.Equal(t, "x509: “JP” certificate is not trusted", ckr.Message)

	ckr = Run([]string{"-H", host, "-p", port, "--no-check-certificate", "-c", "25", "-w", "30"})
	assert.Equal(t, checkers.WARNING, ckr.Status, "should be WARNING")
}
