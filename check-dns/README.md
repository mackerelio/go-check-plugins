# check-dns

## Description

Monitor DNS response.

## Synopsis

Check DNS server status  
If DNS server returns  
```
NOERROR   -> OK  
otherwise -> CRITICAL
```
```
check-dns -H example.com -s 8.8.8.8
```

Check IP-ADDRESS DNS server returns  
If DNS server returns 1.1.1.1 and 2.2.2.2
```
-a 1.1.1.1, 2.2.2.2          -> OK  
-a 1.1.1.1, 2.2.2.2, 3.3.3.3 -> WARNING  
-a 1.1.1.1                   -> WARNING  
-a 1.1.1.1, 3.3.3.3          -> WARNING  
-a 3.3.3.3                   -> CRITICAL  
-a 3.3.3.3, 4.4.4.4, 5.5.5.5 -> CRITICAL  
```
```
check-dns -H example.com -s 8.8.8.8 -a 93.184.216.34
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
  -H, --host=             The name or address you want to query
  -s, --server=           DNS server you want to use for the lookup
  -p, --port=             Port number you want to use (default: 53)
  -q, --querytype=        DNS record query type where TYPE =(A, AAAA, SRV, TXT, MX, ANY) (default: A)
  -c, --queryclass=       DNS record class type where TYPE =(IN, CS, CH, HS, NONE, ANY) (default: IN)
      --norec             Set not recursive mode
  -a, --expected-address= IP-ADDRESS you expect the DNS server to return. If multiple addresses are returned at once, you have to match the whole string of addresses separated with commas
```

## For more information

Please execute `check-dns -h` and you can get command line options.
