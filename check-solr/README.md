# check-solr

## Description

Checks for Apache Solr.

## Synopsis
```
check-solr ping --host=127.0.0.1 --port=8983 --core=CORE
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-solr
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-solr ping --host=127.0.0.1 --port=8983 --core=CORE
```


## Setting for mackerel-agent

If there are no problems in the execution result, add s setting in mackerel-agent.conf .

```
[plugin.checks.check-solr-sample]
command = ["check-solr", "ping", "--host", "127.0.0.1", "--port", "8983", "--core", "CORE"]
```

## Usage
### Subcommands

```
  ping
```

### Options
#### `ping` subcommand

Checks the Apache Solr PING response.

```
  -H, --host= Hostname (default: localhost)
  -p, --port= Port (default: 8983)
  -c, --core= Core
```

## For more information

Please execute `check-solr -h` and you can get command line options.
