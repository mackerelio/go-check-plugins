# check-postgresql

## Description

Checks for PostgreSQL

## Sub Commands

- connection

## check-postgresql connection

Checks the number of PostgreSQL connections.

### Setting

```
[plugin.checks.postgresql_connection]
command = "/path/to/check-postgresql connection --host=127.0.0.1 --port=5432 --user=USER --password=PASSWORD --database=DBNAME --warning=70 --critical=90
```
