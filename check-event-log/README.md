# check-event-log

## Description

Checks a windows event log using a regular expression.

## Setting

```
[plugin.checks.event-log]
command = "/path/to/check-event-log --type=Application --pattern=REGEXP --warning-over=N --critical-over=N"
```

