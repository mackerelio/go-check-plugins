# check-windows-eventlog

## Description

Checks a windows event log using a regular expression.

## Setting

```
[plugin.checks.event-log]
command = "/path/to/check-windows-eventlog --log=LOGTYPE --type=EVENTTYPE --source-pattern=REGEXP --source-exclude=REGEXP --message-pattern=REGEXP --message-exclude=REGEXP --event-id-pattern=RANGE --event-id-exclude=RANGE --warning-over=N --critical-over=N" --fail-first
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

## Tutorial

1. find message matches `foo` but not match `bar`.

    ```
    --message-pattern foo --message-exclude bar
    ```

2. find event which id is 900 or 901.

    ```
    --event-id-pattern 900,901
    ```

3. find event which id is between 900 and 1200, but not 1101.

    ```
    --event-id-pattern 900-1200 --event-id-exclude 1101
    ```
