check-elasticsearch
==========

Check Elasticsearch Health with `/_cluster/health` API.

## Synopsis

```shell
check-elasticsearch [--scheme=<http|https>] [--host=<host>] [--port=<port>]
```

## Setting

```
[plugin.checks.elasticsearch]
command = "/path/to/check-elasticsearch --host=127.0.0.1 --port=9200"
```
