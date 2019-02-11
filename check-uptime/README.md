# check-uptime

## Description

Check uptime seconds.

## Synopsis
```
check-uptime --warning-under=600 --critical-under=120
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-uptime
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-uptime --warning-under=600 --critical-under=120
```


## Setting for mackerel-agent

If there are no problems in the execution result, add s setting in mackerel-agent.conf .

```
[plugin.checks.check-uptime-sample]
command = ["check-uptime", "--warning-under", "600", "--critical-under", "120"]
```

## Usage
### Options

```
      --warn-under=N        (DEPRECATED) Trigger a warning if under the seconds
  -w, --warning-under=N     Trigger a warning if under the seconds
  -c, --critical-under=N    Trigger a critial if under the seconds
      --warn-over=N         (DEPRECATED) Trigger a warning if over the seconds
  -W, --warning-over=N      Trigger a warning if over the seconds
  -C, --critical-over=N     Trigger a critical if over the seconds
```

## For more information

Please execute `check-uptime -h` and you can get command line options.
