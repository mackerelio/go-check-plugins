# check-windows-eventlog

## Description

Checks a windows event log using a regular expression.

## Setting

```
[plugin.checks.event-log]
command = "/path/to/check-windows-eventlog --log=LOGTYPE --type=EVENTTYPE --source-pattern=REGEXP --source-exclude=REGEXP --message-pattern=REGEXP --message-exclude=REGEXP --warning-over=N --critical-over=N" --fail-first
```

## Tutorial

find message matches `foo` but not match `bar`

```
--message-pattern foo --message-exclude bar
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
