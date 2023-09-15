# check-tcp

## Description

This plugin tests TCP connections with the specified host.

## Synopsis
```
check-tcp -H localhost -p 4224 -w 3 -c 5
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-tcp
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-tcp -H localhost -p 4224 -w 3 -c 5
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.check-tcp-sample]
command = ["check-tcp", "-H", "localhost", "-p", "4224", "-w", "3", "-c", "5"]
```

## Usage
### Options

```
      --service=              Service name. e.g. ftp, smtp, pop, imap and so on
  -H, --hostname=             Host name or IP Address
  -p, --port=                 Port number
  -s, --send=                 String to send to the server
  -e, --expect-pattern=       Regexp pattern to expect in server response
  -q, --quit=                 String to send server to initiate a clean close of the connection
  -S, --ssl                   Use SSL for the connection.
  -U, --unix-sock=            Unix Domain Socket
      --no-check-certificate  Do not check certificate
  -t, --timeout=              Seconds before connection times out (default: 10)
  -m, --maxbytes=             Close connection once more than this number of bytes are received
  -d, --delay=                Seconds to wait between sending string and polling for response
  -w, --warning=              Response time to result in warning status (seconds)
  -c, --critical=             Response time to result in critical status (seconds)
  -E, --escape                Can use \n, \r, \t or \ in send or quit string. Must come before send or quit option. By default, nothing added to send, \r\n added to end of quit
  -W, --error-warning         Set the error level to warning when exiting with unexpected error (default: critical). In the case of request succeeded, evaluation result of -c option eval takes priority.
  -C, --expect-closed         Verify that the port/unixsock is closed. If the port/unixsock is closed, OK; if open, follow the ErrWarning flag. This option only verifies the connection.
```

## For more information

Please execute `check-tcp -h` and you can get command line options.

## Other

- [Nagios Plugins - check_tcp](https://www.monitoring-plugins.org/doc/man/check_tcp.html)
