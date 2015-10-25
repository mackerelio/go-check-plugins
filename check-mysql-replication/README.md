# check-mysql-replication
## Description

Checks MySQL replication status and its second behind master.

## Setting

```
[plugin.checks.mysql_replication]
command = "/path/to/check-mysql-replication --host=127.0.0.1 --port=3306 --user=USER --password=PASSWORD --warning=5 --critical=10
```
