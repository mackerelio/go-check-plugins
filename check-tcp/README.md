# check-tcp

## Description

This plugin tests TCP connections with the specified host

## Setting

```
[plugin.checks.tcp]
command = "/path/to/check-tcp -H localhost -p 4224 -w 3 -c 5"
```

## Options

```
    --service=             Service name. e.g. ftp, smtp, pop, imap and so on
-H, --hostname=            Host name or IP Address
-p, --port=                Port number
-s, --send=                String to send to the server
-e, --expect-pattern=      Regexp pattern to expect in server response
-q, --quit=                String to send server to initiate a clean close of the connection
-S, --ssl                  Use SSL for the connection.
    --no-check-certificate Do not check certificate
-U, --unix-sock=           Unix Domain Socket
-t, --timeout=             Seconds before connection times out (default: 10)
-m, --maxbytes=            Close connection once more than this number of bytes are received
-d, --delay=               Seconds to wait between sending string and polling for response
-w, --warning=             Response time to result in warning status (seconds)
-c, --critical=            Response time to result in critical status (seconds)
-E, --escape               Can use \n, \r, \t or \ in send or quit string. Must come before send or quit option. By
                           default, nothing added to send, \r\n added to end of quit
```

## Other

* [Nagios Plugins - check_tcp](https://www.monitoring-plugins.org/doc/man/check_tcp.html)
