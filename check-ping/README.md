# check-ping

## Description
Check ICMP Ping connections with the specified host.

## Synopsis
```
check-ping -H 127.0.0.1 -n 5 -w 100
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-ping
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-ping -H 127.0.0.1 -n 5 -w 100
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.check-ping-sample]
command = ["check-ping", "-H", "127.0.0.1", "-n", "5", "-w", "100"]
```

## Usage
### Options

```
  -H, --host=      check target IP Address
  -n, --count=     sending (and receiving) count ping packets (default: 1)
  -w, --wait-time= wait time, Max RTT(ms) (default: 1000)
      --status-as= Overwrite status=to-status, support multiple comma separetes.
```

## For more information

Please execute `check-ping -h` and you can get command line options.
