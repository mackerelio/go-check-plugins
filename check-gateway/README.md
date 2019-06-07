# check-gateway

## Description

Check ICMP Ping connections with the specified host forcing a specified gateway

Depends on arping-th

## Synopsis
```
check-gateway -H 1.1.1.1 -I eth0 192.168.1.1
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-gateway
go install
```

If you are ok running it as root, you are done!

If you want to run as user, set it setuid root. Probably `setcap` can be enough on Linux, I just haven't had
time to test it.


## Usage
### Options

## For more information

Please execute `check-gateway -h` and you can get command line options.
