package checkmasterha

import (
	"os"
	"strings"
	"testing"

	"github.com/mackerelio/go-check-plugins/check-masterha/lib/mock"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

var subc = subcommand{
	Config:    "/path/to/masterha/db001.conf",
	ConfigDir: "/usr/local/masterha/conf",
	All:       false,
}

var replSubcommand = replChecker{subcommand: subc, SecondsBehindMaster: 0}
var statusSubcommand = statusChecker{subcommand: subc}
var sshSubcommand = sshChecker{subcommand: subc}

func TestMain(m *testing.M) {
	replSubcommand.Executer = &replSubcommand
	statusSubcommand.Executer = &statusSubcommand
	sshSubcommand.Executer = &sshSubcommand
	os.Exit(m.Run())
}

func TestSubcommandExecuteSuccessful(t *testing.T) {
	subc := subcommand{
		Config:    "/path/to/masterha/db001.conf",
		ConfigDir: "/usr/local/masterha/conf",
		All:       false,
		Executer: &mock.Executer{
			CommandName:   "echo",
			CommandArgs:   []string{"hello"},
			CommandResult: "* should set result to here *",
			Status:        checkers.OK,
			ParseResult:   "OK",
		},
	}
	checker := subc.execute("/path/to/masterha/db002.conf")
	assert.Equal(t, "hello --conf /path/to/masterha/db002.conf\n", subc.Executer.(*mock.Executer).CommandResult)
	assert.Equal(t, checkers.OK, checker.Status)
	assert.Equal(t, "OK", checker.Message)
}

func TestSubcommandExecuteUnknown(t *testing.T) {
	subc := subcommand{
		Config:    "/path/to/masterha/db001.conf",
		ConfigDir: "/usr/local/masterha/conf",
		All:       false,
		Executer: &mock.Executer{
			CommandName:   "false",
			CommandArgs:   []string{""},
			CommandResult: "* should set result to here *",
			Status:        checkers.UNKNOWN,
			ParseResult:   "UNKNOWN",
		},
	}
	checker := subc.execute(subc.Config)
	assert.Equal(t, "", subc.Executer.(*mock.Executer).CommandResult)
	assert.Equal(t, checkers.WARNING, checker.Status)
	assert.Equal(t, "UNKNOWN", checker.Message)
}

func TestExtractNonEmptyLines(t *testing.T) {
	lines := strings.Split(`
db001 (pid:1111) is running(0:PING_OK), master:XX.XX.XX.XX
db002 (pid:1111) is running(0:PING_OK), master:XX.XX.XX.XX
db003 (pid:1111) is running(0:PING_OK), master:XX.XX.XX.XX
`, "\n")
	assert.Equal(t, 5, len(lines))
	extracted := extractNonEmptyLines(lines)
	assert.Equal(t, 3, len(extracted))
}

func TestExtractErrorMsg(t *testing.T) {
	body := `Thu Mar 24 14:31:00 2016 - [warn] Global configuration file /etc/masterha_default.cnf not found. Skipping.
Thu Mar 24 14:31:00 2016 - [info] Reading application default configurations from /usr/local/masterha/conf/db001.cnf..
Thu Mar 24 14:31:00 2016 - [info] Reading server configurations from /usr/local/masterha/conf/db001.cnf..
Thu Mar 24 14:31:00 2016 - [info] MHA::MasterMonitor version 0.50.
Thu Mar 24 14:31:01 2016 - [info] Dead Servers:
Thu Mar 24 14:31:01 2016 - [info] Alive Servers:
Thu Mar 24 14:31:01 2016 - [info]   XX.XX.XX.XX(XX.XX.XX.XX:3306)
Thu Mar 24 14:31:01 2016 - [info]   XX.XX.XX.XX(XX.XX.XX.XX:3306)
Thu Mar 24 14:31:01 2016 - [info]   XX.XX.XX.XX(XX.XX.XX.XX:3306)
Thu Mar 24 14:31:01 2016 - [info]   XX.XX.XX.XX(XX.XX.XX.XX:3306)
Thu Mar 24 14:31:01 2016 - [info] Alive Slaves:
Thu Mar 24 14:31:01 2016 - [info]   XX.XX.XX.XX(XX.XX.XX.XX:3306)  Version=5.6.29-log (oldest major version between slaves) log-bin:enabled
Thu Mar 24 14:31:01 2016 - [info]     Replicating from XX.XX.XX.XX(XX.XX.XX.XX:3306)
Thu Mar 24 14:31:01 2016 - [info]     Not candidate for the new Master (no_master is set)
Thu Mar 24 14:31:01 2016 - [info]   XX.XX.XX.XX(XX.XX.XX.XX:3306)  Version=5.6.29-log (oldest major version between slaves) log-bin:enabled
Thu Mar 24 14:31:01 2016 - [info]     Replicating from XX.XX.XX.XX(XX.XX.XX.XX:3306)
Thu Mar 24 14:31:01 2016 - [info]   XX.XX.XX.XX(XX.XX.XX.XX:3306)  Version=5.6.29-log (oldest major version between slaves) log-bin:enabled
Thu Mar 24 14:31:01 2016 - [info]     Replicating from XX.XX.XX.XX(XX.XX.XX.XX:3306)
Thu Mar 24 14:31:01 2016 - [info]   XX.XX.XX.XX(XX.XX.XX.XX:3306)  Version=5.6.29-log (oldest major version between slaves) log-bin:enabled
Thu Mar 24 14:31:01 2016 - [info]     Replicating from XX.XX.XX.XX(XX.XX.XX.XX:3306)
Thu Mar 24 14:31:01 2016 - [error][/usr/lib/perl5/MHA/ServerManager.pm, ln522] FATAL: Replication configuration error. All slaves should replicate from the same master.
Thu Mar 24 14:31:01 2016 - [error][/usr/lib/perl5/MHA/ServerManager.pm, ln1069] MySQL master is not correctly configured. Check master/slave settings
Thu Mar 24 14:31:01 2016 - [error][/usr/lib/perl5/MHA/MasterMonitor.pm, ln316] Error happend on checking configurations.  at /usr/lib/perl5/MHA/MasterMonitor.pm line 243
Thu Mar 24 14:31:01 2016 - [error][/usr/lib/perl5/MHA/MasterMonitor.pm, ln397] Error happened on monitoring servers.
Thu Mar 24 14:31:01 2016 - [info] Got exit code 1 (Not master dead).

MySQL Replication Health is NOT OK!
`
	extracted := extractErrorMsg(body)
	assert.Equal(t, 4, len(strings.Split(extracted, "\n")))
}
