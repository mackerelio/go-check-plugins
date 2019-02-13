# check-file-age

## Description

Monitor file age and size for script monitoring.

## Synopsis
```
check-file-age -w 240 -W 10 -c 600 -C 0 -f /path/to/filename
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-file-age
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-file-age -w 240 -W 10 -c 600 -C 0 -f /path/to/filename
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.fileage-sample]
command = ["check-file-age", "-w", "240", "-W", "10", "-c", "600", "-C", "0", "-f", "/path/to/filename"]
```

## Usage
### Options

```
  -f, --file=           monitor file name
  -w, --warning-age=    warning if more old than (default: 240)
  -W, --warning-size=   warning if file size less than
  -c, --critical-age=   critical if more old than (default: 600)
  -C, --critical-size=  critical if file size less than (default: 0)
  -i, --ignore-missing  skip alert if file doesn't exist
```

## For more information

Please execute `check-file-age -h` and you can get command line options.

## Other

- inspired by [check_file_age.pl](https://github.com/nagios-plugins/nagios-plugins/blob/master/plugins-scripts/check_file_age.pl)
