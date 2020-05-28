package checkmemcached

import (
	"context"
	"net"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func allocUnusedPort() (string, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	_, port, err := net.SplitHostPort(l.Addr().String())
	if err != nil {
		return "", err
	}
	l.Close()
	return port, nil
}

func waitForPort(ctx context.Context, port string) error {
	addr := net.JoinHostPort("127.0.0.1", port)
	d, _ := ctx.Deadline()
	dialer := net.Dialer{
		Deadline: d,
	}
	for ctx.Err() == nil {
		c, err := dialer.Dial("tcp", addr)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		c.Close()
		return nil
	}
	return ctx.Err()
}

func TestMemd(t *testing.T) {
	port, err := allocUnusedPort()
	if err != nil {
		t.Fatal("allocUnusedPort:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "memcached", "-p", port)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	if err := cmd.Start(); err != nil {
		t.Fatal("Start:", err)
	}
	if err := waitForPort(ctx, port); err != nil {
		t.Fatal("waitForPort:", err)
	}

	t.Logf("memcached started on port %s", port)
	defer func() {
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Signal(syscall.SIGTERM)
		}
	}()

	argv := []string{"-p", port, "-k", "test"}
	ckr := run(argv)
	if ckr.Status.String() != "OK" {
		t.Errorf("failed to check memcache:%s", ckr)
	}
	cmd.Process.Signal(syscall.SIGTERM)
	cmd.Wait()
}
