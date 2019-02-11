# check-masterha

## Description

Monitor MasterHA status.

## Synopsis
```
check-masterha status --all
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-masterha
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-masterha status --all
check-masterha repl --conf /path/to/masterha_cluster.cnf
check-masterha ssh --conf /path/to/masterha_cluster.cnf
```


## Setting for mackerel-agent

If there are no problems in the execution result, add s setting in mackerel-agent.conf .

```
[plugin.checks.check-masterha-sample]
command = ["check-masterha", "status", "--all"]
```

## Usage
### Subcommands

```
  repl    check to masterha_check_repl
  ssh     check to masterha_check_ssh
  status  check to masterha_check_status
```

### Options
#### `repl` subcommand

```
      -c, --conf=                  target config file
          --confdir=               config directory (default: /usr/local/masterha/conf)
      -a, --all                    use all config file for target
          --seconds_behind_master= seconds_behind_master option for masterha_check_repl
```

#### `ssh` subcommand

```
      -c, --conf=    target config file
          --confdir= config directory (default: /usr/local/masterha/conf)
      -a, --all      use all config file for target
```

#### `status` subcommand

```
      -c, --conf=    target config file
          --confdir= config directory (default: /usr/local/masterha/conf)
      -a, --all      use all config file for target
```

## For more information

Please execute `check-masterha -h` and you can get command line options.

## References
- [Master HA Manager](https://code.google.com/p/mysql-master-ha/)
