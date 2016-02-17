# check-mysql

## Description

Checks for MySQL

## Sub Commands

- connection
- replication
- uptime

## check-mysql connection

Checks the number of MySQL connections.

### Setting

```
[plugin.checks.mysql_connection]
command = "/path/to/check-mysql connection --host=127.0.0.1 --port=3306 --user=USER --password=PASSWORD --warning=250 --critical=280
```

## check-mysql replication

Checks MySQL replication status and its second behind master.

### Setting

```
[plugin.checks.mysql_replication]
command = "/path/to/check-mysql replication --host=127.0.0.1 --port=3306 --user=USER --password=PASSWORD --warning=5 --critical=10
```
