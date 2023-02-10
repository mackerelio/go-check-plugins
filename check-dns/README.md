# check-dns

## Description

Monitor DNS response.

## Synopsis
```
check-dns -H example.com -s 8.8.8.8
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-dns
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-dns -H example.com -s 8.8.8.8
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.fileage-sample]
command = ["check-dns", "-H", "example.com", "-s", "8.8.8.8"]
```

## Usage
### Options

```
  -H, --host=       The name or address you want to query
  -s, --server=     DNS server you want to use for the lookup
  -p, --port=       Port number you want to use (default: 53)
  -q, --querytype=  DNS record query type where TYPE =(A, AAAA, SRV, TXT, MX, ANY) (default: A)
  -c, --queryclass= DNS record class type where TYPE =(IN, CS, CH, HS, NONE, ANY) (default: IN)
      --norec       Set not recursive mode
```

## For more information

Please execute `check-dns -h` and you can get command line options.
