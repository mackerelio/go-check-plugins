# check-cert-file

## Description
Check expiry for a certification file.


## Synopsis
```
check-cert-file --file=/path/to/cert.pem --warning=30 --critical=14
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-cert-file
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-cert-file --file=/path/to/cert.pem --warning=30 --critical=14
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.cert-file]
command = ["check-cert-file", "--file", "/path/to/cert.pem", "--warning", "30", "--critical", "14"]
```

## Usage
### Options

```
  -f, --file=     cert file name
  -c, --critical= The critical threshold in days before expiry (default: 14)
  -w, --warning=  The threshold in days before expiry (default: 30)
```


## For more information

Please execute `check-cert-file -h` and you can get command line options.
