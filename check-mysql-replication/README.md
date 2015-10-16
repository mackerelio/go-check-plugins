# check-mysql-replication
## Description

Checks MySQL replication status and its second behind master.

## Setting

```
[plugin.checks.mysql-replication]
command = "/path/to/check-mysql-replication -host=127.0.0.1 -port=3306 -username=USER -password=PASSWORD -warn=5 -crit=10
```

