package checkmasterha

import (
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestMakeReplCommandName(t *testing.T) {
	name := replSubcommand.MakeCommandName()
	assert.Equal(t, "masterha_check_repl", name)
}

func TestMakeReplCommandArgs(t *testing.T) {
	args := replSubcommand.MakeCommandArgs()
	assert.Equal(t, 0, len(args))
}

func TestMakeReplCommandArgsWithSecondsBehindMaster(t *testing.T) {
	replSubcommand := replChecker{
		subcommand: subcommand{
			Config:    "/path/to/masterha/db001.conf",
			ConfigDir: "/usr/local/masterha/conf",
			All:       false,
		},
		SecondsBehindMaster: 12345,
	}
	args := replSubcommand.MakeCommandArgs()
	assert.Equal(t, 2, len(args))
	assert.Equal(t, "--seconds_behind_master", args[0])
	assert.Equal(t, "12345", args[1])
}

func TestParseReplSuccess(t *testing.T) {
	body := `Thu Mar 24 14:40:20 2016 - [warn] Global configuration file /etc/masterha_default.cnf not found. Skipping.
Thu Mar 24 14:40:20 2016 - [info] Reading application default configurations from /usr/local/masterha/conf/db001.cnf..
Thu Mar 24 14:40:20 2016 - [info] Reading server configurations from /usr/local/masterha/conf/db001.cnf..
Thu Mar 24 14:40:20 2016 - [info] MHA::MasterMonitor version 0.50.
Thu Mar 24 14:40:21 2016 - [info] Dead Servers:
Thu Mar 24 14:40:21 2016 - [info] Alive Servers:
Thu Mar 24 14:40:21 2016 - [info]   XX.XX.XX.XX(XX.XX.XX.XX:3306)
Thu Mar 24 14:40:21 2016 - [info]   XX.XX.XX.XX(XX.XX.XX.XX:3306)
Thu Mar 24 14:40:21 2016 - [info]   XX.XX.XX.XX(XX.XX.XX.XX:3306)
Thu Mar 24 14:40:21 2016 - [info] Alive Slaves:
Thu Mar 24 14:40:21 2016 - [info]   XX.XX.XX.XX(XX.XX.XX.XX:3306)  Version=5.6.29-log (oldest major version between slaves) log-bin:enabled
Thu Mar 24 14:40:21 2016 - [info]     Replicating from XX.XX.XX.XX(XX.XX.XX.XX:3306)
Thu Mar 24 14:40:21 2016 - [info]   XX.XX.XX.XX(XX.XX.XX.XX:3306)  Version=5.6.29-log (oldest major version between slaves) log-bin:enabled
Thu Mar 24 14:40:21 2016 - [info]     Replicating from XX.XX.XX.XX(XX.XX.XX.XX:3306)
Thu Mar 24 14:40:21 2016 - [info]     Not candidate for the new Master (no_master is set)
Thu Mar 24 14:40:21 2016 - [info] Current Master: XX.XX.XX.XX(XX.XX.XX.XX:3306)
Thu Mar 24 14:40:21 2016 - [info] Checking slave configurations..
Thu Mar 24 14:40:21 2016 - [warn]  relay_log_purge=0 is not set on slave XX.XX.XX.XX(XX.XX.XX.XX:3306).
Thu Mar 24 14:40:21 2016 - [warn]  relay_log_purge=0 is not set on slave XX.XX.XX.XX(XX.XX.XX.XX:3306).
Thu Mar 24 14:40:21 2016 - [info] Checking replication filtering settings..
Thu Mar 24 14:40:21 2016 - [info]  binlog_do_db= , binlog_ignore_db=
Thu Mar 24 14:40:21 2016 - [info]  Replication filtering check ok.
Thu Mar 24 14:40:21 2016 - [info] Starting SSH connection tests..
Thu Mar 24 14:40:22 2016 - [info] All SSH connection tests passed successfully.
Thu Mar 24 14:40:22 2016 - [info] Checking MHA Node version..
Thu Mar 24 14:40:23 2016 - [info]  Version check ok.
Thu Mar 24 14:40:23 2016 - [info] Checking SSH publickey authentication and checking recovery script configurations on the current master..
Thu Mar 24 14:40:23 2016 - [info]   Executing command: save_binary_logs --command=test --start_file=XXXX-bin.XXXXXX --start_pos=1 --binlog_dir=/path/to/lib/mysql,/path/to/log/mysql --output_file=/path/to/log/masterha/db001/save_binary_logs_test --manager_version=0.50
Thu Mar 24 14:40:23 2016 - [info]   Connecting to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
  Creating /path/to/log/masterha/db001 if not exists..    ok.
  Checking output directory is accessible or not..
   ok.
  Binlog found at /path/to/lib/mysql, up to XXXX-bin.XXXXXX
Thu Mar 24 14:40:23 2016 - [info] Master setting check done.
Thu Mar 24 14:40:23 2016 - [info] Checking SSH publickey authentication and checking recovery script configurations on all alive slave servers..
Thu Mar 24 14:40:23 2016 - [info]   Executing command : apply_diff_relay_logs --command=test --slave_user=masterha --slave_host=XX.XX.XX.XX --slave_ip=XX.XX.XX.XX --slave_port=3306 --workdir=/path/to/log/masterha/db001 --target_version=5.6.29-log --manager_version=0.50 --relay_log_info=/path/to/lib/mysql/relay-log.info  --slave_pass=xxx
Thu Mar 24 14:40:23 2016 - [info]   Connecting to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
  Checking slave recovery environment settings..
    Opening /path/to/lib/mysql/relay-log.info XX.XX.XX.XX ok.
    Relay log found at /path/to/lib/mysql, up to XXXXX-relay-bin.XXXXXX
    Temporary relay log file is /path/to/lib/mysql/XXXXX-relay-bin.XXXXXX
    Testing mysql connection and privileges.. done.
    Testing mysqlbinlog output.. done.
    Cleaning up test file(s).. done.
Thu Mar 24 14:40:23 2016 - [info]   Executing command : apply_diff_relay_logs --command=test --slave_user=masterha --slave_host=XX.XX.XX.XX --slave_ip=XX.XX.XX.XX --slave_port=3306 --workdir=/path/to/log/masterha/db001 --target_version=5.6.29-log --manager_version=0.50 --relay_log_info=/path/to/lib/mysql/relay-log.info  --slave_pass=xxx
Thu Mar 24 14:40:23 2016 - [info]   Connecting to mysql@XX.XX.XX.XX(XX.XX.XX.XX)..
  Checking slave recovery environment settings..
    Opening /path/to/lib/mysql/relay-log.info XX.XX.XX.XX ok.
    Relay log found at /path/to/lib/mysql, up to XXXXX-relay-bin.XXXXXX
    Temporary relay log file is /path/to/lib/mysql/XXXXX-relay-bin.XXXXXX
    Testing mysql connection and privileges.. done.
    Testing mysqlbinlog output.. done.
    Cleaning up test file(s).. done.
Thu Mar 24 14:40:24 2016 - [info] Slaves settings check done.
Thu Mar 24 14:40:24 2016 - [info]
XX.XX.XX.XX (current master)
 +--XX.XX.XX.XX
 +--XX.XX.XX.XX

Thu Mar 24 14:40:24 2016 - [info] Checking replication health on XX.XX.XX.XX..
Thu Mar 24 14:40:24 2016 - [info]  ok.
Thu Mar 24 14:40:24 2016 - [info] Checking replication health on XX.XX.XX.XX..
Thu Mar 24 14:40:24 2016 - [info]  ok.
Thu Mar 24 14:40:24 2016 - [info] Checking master_ip_failvoer_script status:
Thu Mar 24 14:40:24 2016 - [info]   /path/to/failover.sh --command=status --ssh_user=root --orig_master_host=XX.XX.XX.XX --orig_master_ip=XX.XX.XX.XX --orig_master_port=3306
Thu Mar 24 14:40:24 2016 - [info]  OK.
Thu Mar 24 14:40:24 2016 - [warn] shutdown_script is not defined.
Thu Mar 24 14:40:24 2016 - [info] Got exit code 0 (Not master dead).

MySQL Replication Health is OK.
`
	status, msg := replSubcommand.Parse(body)
	assert.Equal(t, checkers.OK, status)
	assert.Equal(t, "MySQL Replication Health is OK.", msg)
}

func TestParseReplFailure(t *testing.T) {
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
Thu Mar 24 14:31:01 2016 - [error][/path/to/perl5/MHA/ServerManager.pm, ln522] FATAL: Replication configuration error. All slaves should replicate from the same master.
Thu Mar 24 14:31:01 2016 - [error][/path/to/perl5/MHA/ServerManager.pm, ln1069] MySQL master is not correctly configured. Check master/slave settings
Thu Mar 24 14:31:01 2016 - [error][/path/to/perl5/MHA/MasterMonitor.pm, ln316] Error happend on checking configurations.  at /path/to/perl5/MHA/MasterMonitor.pm line 243
Thu Mar 24 14:31:01 2016 - [error][/path/to/perl5/MHA/MasterMonitor.pm, ln397] Error happened on monitoring servers.
Thu Mar 24 14:31:01 2016 - [info] Got exit code 1 (Not master dead).

MySQL Replication Health is NOT OK!
`
	status, msg := replSubcommand.Parse(body)
	assert.Equal(t, checkers.CRITICAL, status)
	assert.Equal(t, "MySQL Replication Health is NOT OK!", msg)
}
