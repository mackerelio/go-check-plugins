# check-dns

## Description

Monitor DNS response.

## Usage

### Options

```
  -H, --host=            The name or address you want to query
  -s, --server=          DNS server you want to use for the lookup
  -p, --port=            Port number you want to use (default: 53)
  -q, --querytype=       DNS record query type (default: A)
      --norec            Set not recursive mode
  -e, --expected-string= IP-ADDRESS string you expect the DNS server to return. If multiple IP-ADDRESS are returned at once, you have to specify whole string
```

- query class is always IN.
- Punycode is not supported.

### Check DNS server status

If DNS server returns `NOERROR` in status of HEADER, then the checker result becomes `OK`, if not `NOERROR`, then `CRITICAL`

```
check-dns -H a.root-servers.net -s 8.8.8.8
```

### Check string DNS server returns

- The currently supported query types are A, AAAA.

If DNS server returns 1.1.1.1 and 2.2.2.2
```
-e 1.1.1.1 -e 2.2.2.2            -> OK  
-e 1.1.1.1 -e 2.2.2.2 -e 3.3.3.3 -> WARNING  
-e 1.1.1.1                       -> WARNING  
-e 1.1.1.1 -e 3.3.3.3            -> WARNING  
-e 3.3.3.3                       -> CRITICAL  
-e 3.3.3.3 -e 4.4.4.4 -e 5.5.5.5 -> CRITICAL  
```
```
check-dns -H a.root-servers.net -s 8.8.8.8 -e 198.41.0.4
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
check-dns -H a.root-servers.net -s 8.8.8.8
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.dns-sample]
command = ["check-dns", "-H", "a.root-servers.net", "-s", "8.8.8.8"]
```
