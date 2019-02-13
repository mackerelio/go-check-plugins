# check-ssl-cert

## Description

Check the ssl certification's expiry.

## Synopsis
```
check-ssl-cert --host mackerel.io --warning 30 --critical 7
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-ssl-cert
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-ssl-cert --host mackerel.io --warning 30 --critical 7
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.check-ssl-cert-sample]
command = ["check-ssl-cert", "--host", "mackerel.io", "--warning", "30", "--critical", "7"]
```

## Usage
### Options

```
  -H, --host=            Host name
  -p, --port=            Port number (default: 443)
  -w, --warning=days     The warning threshold in days before expiry (default: 30)
  -c, --critical=days    The critical threshold in days before expiry (default: 14)
```

## For more information

Please execute `check-ssl-cert -h` and you can get command line options.
