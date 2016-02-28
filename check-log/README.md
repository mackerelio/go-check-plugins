# check-log

## Description

Checks a log file using a regular expression.

## Setting

```
[plugin.checks.log]
command = "/path/to/check-log --file=/path/to/file --pattern=REGEXP --warning-over=N --critical-over=N"
```

## See Other

* inspired by [sensu-plugins-logs](https://github.com/sensu-plugins/sensu-plugins-logs).
