# check-ftp

## Description

Check for FTP connection.

## Synopsis
```
check-ftp -H ftp.example.com -P 21 -w 3 -c 5 -t 10
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-ftp
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-ftp -H ftp.example.com -w 3 -c 5 -t 10
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.check-ftp-sample]
command = ["check-ftp", "-H", "ftp.example.com", "-P", "21", "-w", "3", "-c", "5", "-t", "10"]
```

## Usage
### Options

```
  -H, --host=                 Hostname (default: localhost)
  -P, --port=                 Port (default: 21)
  -u, --user=                 FTP username (default: anonymous)
  -p, --password=             FTP password (default: anonymous)
  -w, --warning=              Warning threshold (sec)
  -c, --critical=             Critical threshold (sec)
  -t, --timeout=              Timeout (sec) (default: 10)
  -s, --ftps                  Use FTPS
  -i, --implicit-mode         Connects directly using TLS
      --no-check-certificate  Do not check certificate
```

## For more information

Please execute `check-ftp -h` and you can get command line options.
