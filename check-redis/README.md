# check-redis

## Description

Checks for Redis.

## Synopsis
```
check-redis reachable [--host=127.0.0.1] [--port=6379] [--timeout=5] [--socket=<unix socket>]
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-redis
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-redis reachable --host=127.0.0.1 --port=6379 --timeout=5 --socket=<unix socket>
check-redis slave --host=127.0.0.1 --port=6379 --timeout=5 --socket=<unix socket>
```


## Setting for mackerel-agent

If there are no problems in the execution result, add s setting in mackerel-agent.conf .

```
[plugin.checks.check-redis-sample]
command = ["check-redis", "reachable", "--host", "127.0.0.1", "--port", "6379", "--timeout", "5", "--socket", "<unix socket>"]
```

## Usage
### Subcommands

```
  reachable
  replication
  slave
```

### Options
#### `reachable` subcommand

Checks if Redis is reachable.

```
  -H, --host=    Hostname (default: localhost)
  -s, --socket=  Server socket
  -p, --port=    Port (default: 6379)
  -t, --timeout= Dial Timeout in sec (default: 5)
```

#### `replication` subcommand

Check if Redis's replication is working properly.

```
  -H, --host=        Hostname (default: localhost)
  -s, --socket=      Server socket
  -p, --port=        Port (default: 6379)
  -t, --timeout=     Dial Timeout in sec (default: 5)
      --skip-master  return ok if redis role is master
```

#### **【DEPRECATED】** `slave` subcommand

Checks Redis slave status. This subcommand is deprecated. Please use the `replication` subcommand.

```
  -H, --host=    Hostname (default: localhost)
  -s, --socket=  Server socket
  -p, --port=    Port (default: 6379)
  -t, --timeout= Dial Timeout in sec (default: 5)
```

## For more information

Please execute `check-redis -h` and you can get command line options.
