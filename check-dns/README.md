# check-dns

## Description

Monitor DNS response.

## Usage

### Options

```
  -H, --host=            The name or address you want to query
  -s, --server=          DNS server you want to use for the lookup
  -p, --port=            Port number you want to use (default: 53)
  -q, --querytype=       DNS record query type where TYPE =(A, AAAA, TXT, MX, CNAME) (default: A)
  -c, --queryclass=      DNS record class type where TYPE =(IN, CS, CH, HS, NONE, ANY) (default: IN)
      --norec            Set not recursive mode
  -e, --expected-string= The string you expect the DNS server to return. If multiple responses are returned at once, you have to specify whole string
```

- The currently supported query types are A, AAAA, TXT, MX, CNAME.
- Punycode is not supported.

### Check DNS server status

If DNS server returns `NOERROR` in status of HEADER, then the checker result becomes `OK`, if not `NOERROR`, then `CRITICAL`

```
check-dns -H example.com -s 8.8.8.8
```

### Check string DNS server returns

If DNS server returns 1.1.1.1 and 2.2.2.2
```
-a 1.1.1.1 -a 2.2.2.2            -> OK  
-a 1.1.1.1 -a 2.2.2.2 -a 3.3.3.3 -> WARNING  
-a 1.1.1.1                       -> WARNING  
-a 1.1.1.1 -a 3.3.3.3            -> WARNING  
-a 3.3.3.3                       -> CRITICAL  
-a 3.3.3.3 -a 4.4.4.4 -a 5.5.5.5 -> CRITICAL  
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
[plugin.checks.dns-sample]
command = ["check-dns", "-H", "example.com", "-s", "8.8.8.8"]
```

## NOTICES AND INFORMATION

This plugin incorporates material from third parties.

### miekg/dns

**source**: https://github.com/miekg/dns

```
BSD 3-Clause License

Copyright (c) 2009, The Go Authors. Extensions copyright (c) 2011, Miek Gieben. 
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
```
