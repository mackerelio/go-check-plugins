# check-mailq

## Description

Monitor mail queue count.

## Synopsis
```
check-mailq -w 100 -c 200 -M postfix
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-mailq
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-mailq -w 100 -c 200 -M postfix
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.check-mailq-sample]
command = ["check-mailq", "-w", "100", "-c", "200", "-M", "postfix"]
```

## Usage
### Options

```
  -w, --warning=  number of messages in queue to generate warning (default: 100)
  -c, --critical= number of messages in queue to generate critical alert ( w < c ) (default: 200)
  -M, --mta=      target mta (default: postfix)
```

## For more information

Please execute `check-mailq -h` and you can get command line options.

## other
- inspired by [check_mailq.pl](https://github.com/nagios-plugins/nagios-plugins/blob/master/plugins-scripts/check_mailq.pl)
