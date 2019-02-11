# check-ldap

## Description

Check LDAP response time.

## Synopsis
```
check-ldap -b 'dc=jp' -w 5 -c 10
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-ldap
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-ldap -b 'dc=jp' -w 5 -c 10
```


## Setting for mackerel-agent

If there are no problems in the execution result, add s setting in mackerel-agent.conf .

```
[plugin.checks.check-ldap-sample]
command = ["check-ldap", "-b", "dc=jp", "-w", "5", "-c", "10"]
```

## Usage
### Options

```
  -w, --warning=  Response time to result in warning status (seconds)
  -c, --critical= Response time to result in critical status (seconds)
  -H, --host=     Hostname (default: localhost)
  -p, --port=     Port number (default: 389)
  -b, --base=     LDAP base
  -a, --attr=     LDAP attribute to search (default: (objectclass=*))
  -D, --bind=     LDAP bind DN
  -P, --password= LDAP password
```

## For more information

Please execute `check-ldap -h` and you can get command line options.
