# check-disk

## Description
Check free space of disk.

## Synopsis
```
check-disk --warning=10 --critical=5 --path=/
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-disk
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-disk --warning=10 --critical=5 --path=/
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.disk-sample]
command = ["check-disk", "--warning", "10", "--critical", "5", "--path", "/"]
```

## Usage
### Options

```
  -w, --warning=N, N%                  Exit with WARNING status if less than N units or N% of disk are free
  -c, --critical=N, N%                 Exit with CRITICAL status if less than N units or N% of disk are free
  -W, --iwarning=N%                    Exit with WARNING status if less than PERCENT of inode space is free
  -K, --icritical=N%                   Exit with CRITICAL status if less than PERCENT of inode space is free
  -p, --path=PATH                      Mount point or block device as emitted by the mount(8) command (may be repeated)
  -x, --exclude-device=EXCLUDE PATH    Ignore device (may be repeated; only works if -p unspecified)
  -A, --all                            Explicitly select all paths.
  -X, --exclude-type=TYPE              Ignore all filesystems of indicated type (may be repeated)
  -N, --include-type=TYPE              Check only filesystems of indicated type (may be repeated)
  -u, --units=STRING                   Choose bytes, kB, MB, GB, TB (default: MB)
```

## For more information
Please refer to the following.

- Read Mackerel Docs; [Disk monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/check/disk)
- Execute `check-disk -h` and you can get command line options.
