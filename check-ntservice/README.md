# check-ntservice

## Description
Checks Windows NT Service is stopped.


## Synopsis
```
check-ntservice --service-name=SERVICE_NAME
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-ntservice
go install
```

Or when you installing the mackerel-agent msi package, this plug-in is included in the installation folder. About installing mackerel-agent in Windows, see [Installing mackerel-agent on Windows - Mackerel Docs](https://mackerel.io/docs/entry/howto/install-agent/msi).


Next, you can execute this program :-)

```
check-ntservice --service-name=SERVICE_NAME
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.check-ntservice-sample]
command = ["check-ntservice", "--service-name", "SERVICE_NAME"]
```

## Usage
### Options

```
  -s, --service-name=    service name
  -E, --exclude-service= service name to exclude from matching
  -l, --list-service     list service
```


## For more information

Please execute `check-ntservice -h` and you can get command line options.
