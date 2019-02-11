# check-file-size

## Description

Check file size in specified directory.

## Synopsis
```
check-file-size -b /fluentd/buffer_dir/ -w 5M -c 10M
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-file-size
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-file-size -b /fluentd/buffer_dir/ -w 5M -c 10M
```


## Setting for mackerel-agent

If there are no problems in the execution result, add s setting in mackerel-agent.conf .

```
[plugin.checks.filesize-sample]
command = ["check-file-size", "-b", "/fluentd/buffer_dir/", "-w", "5M", "-c", "10M"]
```

## Usage
### Options

```
  -b, --base=     base directory
  -w, --warning=  warning if the size is over (default: 1K)
  -c, --critical= critical if the size is over (default: 1K)
  -d, --depth=    max depth of the directory from base directory (default: 1)
```

## For more information

Please execute `check-file-size -h` and you can get command line options.
