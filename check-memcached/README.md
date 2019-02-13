# check-memcached

## Description
Check memcached by set and get specified key.


## Synopsis
```
check-memcached -H localhost -p 11211 -t 3 -k KeyForTest
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-memcached
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-memcached -H localhost -p 11211 -t 3 -k KeyForTest
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.cert-memcached-sample]
command = ["check-memcached", "-H", "<Hostname>", "-p", "<Port>", "-t", "<DialTimeout>", "-k", "<KeyForTest>"]
```

## Usage
### Options

```
  -H, --host=    Hostname (default: localhost)
  -p, --port=    Port (default: 11211)
  -t, --timeout= Dial Timeout in sec (default: 3)
  -k, --key=     Cache key used within set and get test
```


## For more information

Please execute `check-memcached -h` and you can get command line options.
