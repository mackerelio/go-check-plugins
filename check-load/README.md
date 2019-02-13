# check-load

## Description

Check system load average.

## Synopsis
```
check-load -w 4,3,2 -c 3,2,1
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-load
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-load -w 4,3,2 -c 3,2,1
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.check-load-sample]
command = ["check-load", "-w", "4,3,2", "-c", "3,2,1"]
```

## Usage
### Options

```
  -w, --warning=WL1,WL5,WL15     Warning threshold for loadavg1,5,15
  -c, --critical=CL1,CL5,CL15    Critical threshold for loadavg1,5,15
  -r, --percpu                   Divide the load averages by cpu count
```

## For more information

Please execute `check-load -h` and you can get command line options.

## References
- [check_load](https://github.com/nagios-plugins/nagios-plugins/blob/master/plugins/check_load.c)