package checkmasterha

import (
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestMakeSshCommandName(t *testing.T) {
	name := sshSubcommand.MakeCommandName()
	assert.Equal(t, "masterha_check_ssh", name)
}

func TestMakeSshCommandArgs(t *testing.T) {
	args := sshSubcommand.MakeCommandArgs()
	assert.Equal(t, 0, len(args))
}

func TestParseSshSuccess(t *testing.T) {
	body := `Thu Mar 24 15:42:56 2016 - [warn] Global configuration file /etc/masterha_default.cnf not found. Skipping.
Thu Mar 24 15:42:56 2016 - [info] Reading application default configurations from /usr/local/masterha/conf/db001.cnf..
Thu Mar 24 15:42:56 2016 - [info] Reading server configurations from /usr/local/masterha/conf/db001.cnf..
Thu Mar 24 15:42:56 2016 - [info] Starting SSH connection tests..
Thu Mar 24 15:42:57 2016 - [debug]
Thu Mar 24 15:42:56 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:56 2016 - [debug]   ok.
Thu Mar 24 15:42:56 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:56 2016 - [debug]   ok.
Thu Mar 24 15:42:56 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:56 2016 - [debug]   ok.
Thu Mar 24 15:42:57 2016 - [debug]
Thu Mar 24 15:42:56 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:56 2016 - [debug]   ok.
Thu Mar 24 15:42:56 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:57 2016 - [debug]   ok.
Thu Mar 24 15:42:57 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:57 2016 - [debug]   ok.
Thu Mar 24 15:42:58 2016 - [debug]
Thu Mar 24 15:42:57 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:57 2016 - [debug]   ok.
Thu Mar 24 15:42:57 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:57 2016 - [debug]   ok.
Thu Mar 24 15:42:57 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:57 2016 - [debug]   ok.
Thu Mar 24 15:42:58 2016 - [debug]
Thu Mar 24 15:42:57 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:57 2016 - [debug]   ok.
Thu Mar 24 15:42:57 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:58 2016 - [debug]   ok.
Thu Mar 24 15:42:58 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:58 2016 - [debug]   ok.
Thu Mar 24 15:42:58 2016 - [info] All SSH connection tests passed successfully.
`
	status, msg := sshSubcommand.Parse(body)
	assert.Equal(t, checkers.OK, status)
	assert.Equal(t, "Thu Mar 24 15:42:58 2016 - [info] All SSH connection tests passed successfully.", msg)
}

func TestParseSshFailure(t *testing.T) {
	body := `Thu Mar 24 15:42:56 2016 - [warn] Global configuration file /etc/masterha_default.cnf not found. Skipping.
Thu Mar 24 15:42:56 2016 - [info] Reading application default configurations from /usr/local/masterha/conf/db001.cnf..
Thu Mar 24 15:42:56 2016 - [info] Reading server configurations from /usr/local/masterha/conf/db001.cnf..
Thu Mar 24 15:42:56 2016 - [info] Starting SSH connection tests..
Thu Mar 24 15:42:57 2016 - [debug]
Thu Mar 24 15:42:56 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:56 2016 - [debug]   ok.
Thu Mar 24 15:42:56 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:56 2016 - [debug]   ok.
Thu Mar 24 15:42:56 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:56 2016 - [debug]   ok.
Thu Mar 24 15:42:57 2016 - [debug]
Thu Mar 24 15:42:56 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:56 2016 - [debug]   ok.
Thu Mar 24 15:42:56 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:57 2016 - [debug]   ok.
Thu Mar 24 15:42:57 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:57 2016 - [debug]   ok.
Thu Mar 24 15:42:58 2016 - [debug]
Thu Mar 24 15:42:57 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:57 2016 - [debug]   ok.
Thu Mar 24 15:42:57 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:57 2016 - [debug]   ok.
Thu Mar 24 15:42:57 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:57 2016 - [debug]   ok.
Thu Mar 24 15:42:58 2016 - [debug]
Thu Mar 24 15:42:57 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:57 2016 - [debug]   ok.
Thu Mar 24 15:42:57 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:58 2016 - [debug]   ok.
Thu Mar 24 15:42:58 2016 - [debug]  Connecting via SSH from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
Thu Mar 24 15:42:58 2016 - [error]  SSH connection from mysql@XX.XX.XX.XX(XX.XX.XX.XX) to mysql@XX.XX.XX.XX(XX.XX.XX.XX) failed!
SSH Configuration Check Failed!
`
	status, msg := sshSubcommand.Parse(body)
	assert.Equal(t, checkers.CRITICAL, status)
	assert.Equal(t, "SSH Configuration Check Failed!", msg)
}

