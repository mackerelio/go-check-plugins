package checkmemcached

import (
	"fmt"
	"log"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/lestrrat/go-tcptest"
)

func TestMemd(t *testing.T) {
	var cmd *exec.Cmd
	memd := func(t *tcptest.TCPTest) {
		cmd = exec.Command("memcached", "-p", fmt.Sprintf("%d", t.Port()))
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}
		cmd.Run()
	}

	server, err := tcptest.Start2(memd, 30*time.Second)
	if err != nil {
		log.Fatalf("Failed to start memcached: %s", err)
	}
	t.Logf("memcached started on port %d", server.Port())
	defer func() {
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Signal(syscall.SIGTERM)
		}
	}()

	argv := []string{"-p", fmt.Sprintf("%d", server.Port()), "-k", "test"}
	ckr := run(argv)
	if ckr.Status.String() != "OK" {
		t.Errorf("faild to check memcache:%s", ckr)
	}
	cmd.Process.Signal(syscall.SIGTERM)
	server.Wait()
}
