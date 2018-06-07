# check-masterha

## Description

Monitor MasterHA status.

## Usage

```
$ check-masterha --conf /path/to/masterha

```

## Setting

```
[plugin.checks.masterha-status]
command = "/path/to/check-masterha status --all"

[plugin.checks.masterha-repl]
command = "/path/to/check-masterha repl --conf /path/to/masterha_cluster.cnf"

[plugin.checks.masterha-ssh]
command = "/path/to/check-masterha ssh --conf /path/to/masterha_cluster.cnf"
```

## Other

* [Master HA Manager](https://code.google.com/p/mysql-master-ha/)
