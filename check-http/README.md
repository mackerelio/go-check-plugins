# check-http

## Description

Check about http connection.

## Synopsis
```
check-http -u http://example.com
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-http
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-http -u http://example.com
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.check-http-sample]
command = ["check-http", "-u", "http://example.com"]
```

## Usage
### Options

```
  -u, --url=                                  A URL to connect to
  -s, --status=                               mapping of HTTP status
      --no-check-certificate                  Do not check certificate
  -i, --source-ip=                            source IP address
  -H=                                         HTTP request headers
  -p, --pattern=                              Expected pattern in the content
      --max-redirects=                        Maximum number of redirects followed (default: 10)
      --connect-to=HOST1:PORT1:HOST2:PORT2    Request to HOST2:PORT2 instead of HOST1:PORT1
  -x, --proxy=[PROTOCOL://]HOST[:PORT]        Use the specified proxy. PROTOCOL's default is http, and PORT's default is 1080.
```


To override status
```shell
check-http -s 404=ok -u http://example.com
check-http -s 200-404=ok -u http://example.com
```

To change request destination
```shell
check-http --connect-to=example.com:443:127.0.0.1:8080 https://example.com # will request to 127.0.0.1:8000 but AS example.com:443
check-http --connect-to=:443:127.0.0.1:8080 https://example.com # empty host1 matches ANY host
check-http --connect-to=example.com::127.0.0.1:8080 https://example.com # empty port1 matches ANY port
check-http --connect-to=localhost:443::8080 https://localhost # empty host2 means unchanged, therefore will request to localhost:8080 AS localhost:443
check-http --connect-to=example.com:443:127.0.0.1: https://example.com # empty port2 means unchanged, therefore will request to 127.0.0.1:443
```

To request via proxy (http/https/socks5)
```shell
check-http --proxy=http://localhost:8080 -u http://example.com # request via http://localhost:8080
HTTP_PROXY=http://localhost:8080 check-http  -u http://example.com # Same. you can set proxy via environment variable
check-http --proxy=http://user:pass@localhost:8080 -u http://example.com # basic authentication is also supported
```
## For more information

Please execute `check-http -h` and you can get command line options.
