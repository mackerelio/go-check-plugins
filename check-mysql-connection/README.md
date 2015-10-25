# check-mysql-connection
## Description

Checks the number of MySQL connections.

## Setting

```
[plugin.checks.mysql-connection]
command = "/path/to/check-mysql-connection --host=127.0.0.1 --port=3306 --user=USER --password=PASSWORD --warning=250 --critical=280
```

