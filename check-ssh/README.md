# check-ssh

## Description

Monitor SSHD status.

## Synopsis
```
check-ssh -w 1 -c 3
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-ssh
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-ssh -w 1 -c 3
```


## Setting for mackerel-agent

If there are no problems in the execution result, add s setting in mackerel-agent.conf .

```
[plugin.checks.check-ssh-sample]
command = ["check-ssh", "-w", "1", "-c", "3"]
```

## Usage
### Options

```
  -H, --hostname=   Host name or IP Address (default: localhost)
  -P, --port=       Port number (default: 22)
  -t, --timeout=    Seconds before connection times out (default: 30)
  -w, --warning=    Response time to result in warning status (seconds)
  -c, --critical=   Response time to result in critical status (seconds)
  -u, --user=       Login user name [$USER]
  -p, --password=   Login password [$LOGIN_PASSWORD]
  -i, --identity=   Identity file (ssh private key)
      --passphrase= Identity passphrase [$CHECK_SSH_IDENTITY_PASSPHRASE]
```

## For more information

Please execute `check-ssh -h` and you can get command line options.
