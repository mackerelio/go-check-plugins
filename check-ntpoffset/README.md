# check-ntpoffset

## Description
Check ntp offset.


## Synopsis
```
check-ntpoffset -w=50 -c=100
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-ntpoffset
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-ntpoffset -w=50 -c=100
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.cert-file]
command = ["check-ntpoffset", "-w", "50", "-c", "100"]
```

## Usage
### Options

```
  -c, --critical=     Critical threshold of ntp offset(ms) (default: 100)
  -w, --warning=      Warning threshold of ntp offset(ms) (default: 50)
  -s, --ntp-servers=  Use specified NTP Servers(plural servers can be set separated by ,). When set plural servers, use first
                      response. If not set, use local command just like ntpd/chronyd.
  -t, --ntp-timeout=  Timeout of NTP Server Querying(in seconds). (default: 15)
  -S, --check-stratum Check stratum and fail if the machine is not synchronized.
```


## For more information

Please execute `check-ntpoffset -h` and you can get command line options.
