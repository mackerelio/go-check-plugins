# check-procs

## Description

Checks a log file for a regular expression.

## Setting

```
[plugin.checks.log]
command = "/path/to/check-log --pattern=REGEXP --warning-over=N --critical-over=N"
```

## See Other

* inspired by [sensu-plugins-logs](https://github.com/sensu-plugins/sensu-plugins-logs).
