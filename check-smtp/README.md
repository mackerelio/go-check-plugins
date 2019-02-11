# check-smtp

## Description

Check for SMTP connection.

## Synopsis
```
check-smtp -H smtp.example.com -p 25 -w 3 -c 5 -t 10
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-smtp
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-smtp -H smtp.example.com -p 25 -w 3 -c 5 -t 10
```


## Setting for mackerel-agent

If there are no problems in the execution result, add s setting in mackerel-agent.conf .

```
[plugin.checks.check-smtp-sample]
command = ["check-smtp", "-H", "smtp.example.com", "-p", "25", "-w", "3", "-c", "5", "-t", "10"]
```

## Usage
### Options

```
  -H, --host=         Hostname (default: localhost)
  -p, --port=         Port (default: 25)
  -F, --fqdn=         FQDN used for HELO
  -s, --smtps         Use SMTP over TLS
  -S, --starttls      Use STARTTLS
  -A, --authmech=     SMTP AUTH Authentication Mechanisms (only PLAIN supported)
  -U, --authuser=     SMTP AUTH username
  -P, --authpassword= SMTP AUTH password
  -w, --warning=      Warning threshold (sec)
  -c, --critical=     Critical threshold (sec)
  -t, --timeout=      Timeout (sec) (default: 10)
```

## For more information

Please execute `check-smtp -h` and you can get command line options.
