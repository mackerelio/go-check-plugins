# check-jmx-jolokia

## Description

Checks for JMX value using jolokia.

## Synopsis
```
check-jmx-jolokia -H 127.0.0.1 -p 8778 -m java.lang:type=OperatingSystem -a ProcessCpuLoad -w 10 -c 20
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-jmx-jolokia
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-jmx-jolokia -H 127.0.0.1 -p 8778 -m java.lang:type=OperatingSystem -a ProcessCpuLoad -w 10 -c 20
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.check-jmx-jolokia-sample]
command = ["check-jmx-jolokia", "-H", "127.0.0.1", "-p", "8778", "-m", "java.lang:type=OperatingSystem", "-a", "ProcessCpuLoad", "-w", "10", "-c", "20"]
```

## Usage
### Options

```
  -H, --host=       Host name or IP Address
  -p, --port=       Port (default: 8778)
  -t, --timeout=    Seconds before connection times out (default: 10)
  -m, --mbean=      MBean
  -a, --attribute=  Attribute
  -i, --inner-path= InnerPath
  -k, --key=        Key (default: value)
  -w, --warning=    Trigger a warning if over a number
  -c, --critical=   Trigger a critical if over a number
```

## For more information

Please execute `check-jmx-jolokia -h` and you can get command line options.
