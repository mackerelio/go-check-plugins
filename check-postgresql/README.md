# check-postgresql

## Description
Checks for PostgreSQL.


## Synopsis
```
check-postgresql connection --host=127.0.0.1 --port=5432 --user=USER --password=PASSWORD --database=DBNAME --warning=70 --critical=90
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-postgresql
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-postgresql connection --host=127.0.0.1 --port=5432 --user=USER --password=PASSWORD --database=DBNAME --warning=70 --critical=90
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.check-postgresql-sample]
command = ["check-postgresql", "connection", "--host", "127.0.0.1", "--port", "5432", "--user", "USER", "--password", "PASSWORD", "--database", "DBNAME", "--warning", "70", "--critical", "90"]
```

## Usage
### Subcommands

```
  connection
```

### Options
#### `connection` subcommand

Checks the number of PostgreSQL connections.

```
  -H, --host=        Hostname (default: localhost)
  -p, --port=        Port (default: 5432)
  -u, --user=        Username (default: postgres)
  -P, --password=    Password [$PGPASSWORD]
  -d, --database=    DBname
  -s, --sslmode=     SSLmode (default: disable)
      --sslrootcert= The root certificate used for SSL certificate verification.
  -t, --timeout=     Maximum wait for connection, in seconds. (default: 5)
  -w, --warning=     warning if the number of connection is over (default: 70)
  -c, --critical=    critical if the number of connection is over (default: 90)

```

## For more information

Please execute `check-postgresql -h` and you can get command line options.
