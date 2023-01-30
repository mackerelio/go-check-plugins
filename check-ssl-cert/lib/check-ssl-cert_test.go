package checksslcert

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

const (
	cert_pem   = "cert.pem"
	cert_key   = "key.pem"
	client_crt = "client.crt"
	client_key = "client.key"
)

// generate "cert.pem", "key.pem", "client.crt", "client.key" in tmpDir
func prepareCertification(t *testing.T) (string, error) {
	tmpDir := t.TempDir()
	cmd := exec.Command(
		"go",
		"run",
		filepath.FromSlash(filepath.Join(runtime.GOROOT(), "/src/crypto/tls/generate_cert.go")),
		"-host",
		"127.0.0.1",
		"-duration",
		"720h0m0s", // 30*24*time.Hour
	)
	cmd.Dir = tmpDir
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	os.Rename(filepath.Join(tmpDir, "cert.pem"), filepath.Join(tmpDir, "client.crt"))
	os.Rename(filepath.Join(tmpDir, "key.pem"), filepath.Join(tmpDir, "client.key"))

	cmd2 := exec.Command(
		"go",
		"run",
		filepath.FromSlash(filepath.Join(runtime.GOROOT(), "/src/crypto/tls/generate_cert.go")),
		"-host",
		"127.0.0.1",
		"-duration",
		"720h0m0s", // 30*24*time.Hour
	)
	cmd2.Dir = tmpDir
	err = cmd2.Run()
	if err != nil {
		return "", err
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

	cert, err := tls.LoadX509KeyPair(filepath.Join(tmpDir, cert_pem), filepath.Join(tmpDir, cert_key))
	if err != nil {
		t.Errorf("Failed to LoadX509KeyPair: %s", err)
	}
	ts.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}

	ts.StartTLS()
	defer ts.Close()

	u, _ := url.Parse(ts.URL)
	host, port, _ := net.SplitHostPort(u.Host)

	ckr := Run([]string{
		"-H", host,
		"-p", port,
		"--ca-file", filepath.Join(tmpDir, cert_pem),
		"-c", "25",
		"-w", "30"})
	t.Logf("%s \n", ckr.String())
	assert.Equal(t, checkers.WARNING, ckr.Status, "should be WARNING")
}

func TestClientCertification(t *testing.T) {
	tmpDir, err := prepareCertification(t)
	if err != nil {
		t.Errorf("Failed to prepare: %s", err)
	}

	mux := http.NewServeMux()
	ts := httptest.NewUnstartedServer(mux)

	caCertPEM, err := os.ReadFile(filepath.Join(tmpDir, cert_pem))
	if err != nil {
		t.Errorf("Failed to read file: %s", err)
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(caCertPEM)
	if !ok {
		panic("failed to parse root certificate")
	}

	cert, err := tls.LoadX509KeyPair(filepath.Join(tmpDir, cert_pem), filepath.Join(tmpDir, cert_key))
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
		"--ca-file", filepath.Join(tmpDir, "cert.pem"),
		"--cert-file", filepath.Join(tmpDir, "client.crt"),
		"--key-file", filepath.Join(tmpDir, "client.key"),
		"-c", "25",
		"-w", "30"})
	t.Logf("%s \n", ckr.String())
	assert.Equal(t, checkers.WARNING, ckr.Status, "should be WARNING")
}

func TestNoCheckCertificate(t *testing.T) {
	tmpDir, err := prepareCertification(t)
	if err != nil {
		t.Errorf("Failed to prepare: %s", err)
	}

	mux := http.NewServeMux()
	ts := httptest.NewUnstartedServer(mux)

	cert, err := tls.LoadX509KeyPair(filepath.Join(tmpDir, cert_pem), filepath.Join(tmpDir, cert_key))
	if err != nil {
		t.Errorf("Failed to LoadX509KeyPair: %s", err)
	}
	ts.TLS = &tls.Config{Certificates: []tls.Certificate{cert}}

	ts.StartTLS()

	u, _ := url.Parse(ts.URL)
	host, port, _ := net.SplitHostPort(u.Host)

	ckr := Run([]string{"-H", host, "-p", port, "-c", "25", "-w", "30"})
	t.Logf("%s \n", ckr.String())
	assert.Equal(t, checkers.CRITICAL, ckr.Status, "should be CRITICAL")
	assert.Contains(t, ckr.Message, "x509:", "should an error occur regarding x509")

	ckr = Run([]string{"-H", host, "-p", port, "--no-check-certificate", "-c", "25", "-w", "30"})
	t.Logf("%s \n", ckr.String())
	assert.Equal(t, checkers.WARNING, ckr.Status, "should be WARNING")
}
