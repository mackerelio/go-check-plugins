# check-elasticsearch

## Description
Check Elasticsearch Health with `/_cluster/health` API.

## Synopsis
```
check-elasticsearch [--scheme=<http|https>] [--host=<host>] [--port=<port>]
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-elasticsearch
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-elasticsearch --host=127.0.0.1 --port=9200
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.es-sample]
command = ["check-elasticsearch", "--host", "127.0.0.1", "--port", "9200"]
```

## Usage
### Options

```
  -s, --scheme= Elasticsearch scheme (default: http)
  -H, --host=   Elasticsearch host (default: localhost)
  -p, --port=   Elasticsearch port (default: 9200)
```

## For more information

Please execute `check-elasticsearch -h` and you can get command line options.
