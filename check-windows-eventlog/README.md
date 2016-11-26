# check-windows-eventlog

## Description

Checks a windows event log using a regular expression.

## Setting

```
[plugin.checks.event-log]
command = "/path/to/check-windows-eventlog --log=LOGTYPE --type=EVENTTYPE --source-pattern=REGEXP --message-pattern=REGEXP --warning-over=N --critical-over=N" --fail-first
```

### LOGTYPE

* Application
* Security
* System

### EVENTTYPE

* Success
* Error
* Audit Failure
* Audit Success
* Information
* Warning
