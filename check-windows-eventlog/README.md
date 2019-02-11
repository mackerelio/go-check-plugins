# check-windows-eventlog

## Description

Checks a windows event log using a regular expression.

## Synopsis
```
check-windows-eventlog --log=LOGTYPE --type=EVENTTYPE --source-pattern=REGEXP --source-exclude=REGEXP --message-pattern=REGEXP --message-exclude=REGEXP --event-id-pattern=RANGE --event-id-exclude=RANGE --warning-over=N --critical-over=N --fail-first
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-windows-eventlog
go install
```

Or when you installing the mackerel-agent msi package, this plug-in is included in the installation folder. About installing mackerel-agent in Windows, see [Installing mackerel-agent on Windows - Mackerel Docs](https://mackerel.io/docs/entry/howto/install-agent/msi).


Next, you can execute this program :-)

```
check-windows-eventlog --log=LOGTYPE --type=EVENTTYPE --source-pattern=REGEXP --source-exclude=REGEXP --message-pattern=REGEXP --message-exclude=REGEXP --event-id-pattern=RANGE --event-id-exclude=RANGE --warning-over=N --critical-over=N --fail-first
```


## Setting for mackerel-agent

If there are no problems in the execution result, add s setting in mackerel-agent.conf .

```
[plugin.checks.check-windows-eventlog-sample]
command = ["check-windows-eventlog", "--log", "LOGTYPE", "--type", "EVENTTYPE", "--source-pattern", "REGEXP", "--source-exclude", "REGEXP", "--message-pattern", "REGEXP", "--message-exclude", "REGEXP", "--event-id-pattern", "RANGE", "--event-id-exclude", "RANGE", "--warning-over", "N", "--critical-over", "N", "--fail-first"]
```

## Usage
### Options

```
      --log               Event Names (comma separated)
      --type              Event Types (comma separated)
      --source-pattern    Event Source (regexp pattern)
      --source-exclude    Event Source excluded (regexp pattern)
      --message-pattern   Message Pattern (regexp pattern)
      --message-exclude   Message Pattern excluded (regexp pattern)
      --event-id-pattern  Event IDs acceptable (separated by comma, or range)
      --event-id-exclude  Event IDs ignorable (separated by comma, or range)
  -w, --warning-over      Trigger a warning if matched lines is over a number
  -c, --critical-over     Trigger a critical if matched lines is over a number
  -r, --return            Return matched line
  -s, --state-dir         Dir to keep state files under
      --no-state          Don't use state file and read whole logs
      --fail-first        Count errors on first seek
      --verbose           Verbose output
```

#### LOGTYPE

- Application
- Security
- System

#### EVENTTYPE

- Error
- Audit Failure
- Warning

The following EVENTTYPE can not be detected as an alert.

- Success
- Audit Success
- Information

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

## For more information

Please execute `check-windows-eventlog -h` and you can get command line options.
