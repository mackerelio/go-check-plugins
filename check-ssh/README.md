# check-ssh

## Description

Monitor SSHD status.

## Usage

```
$ check-ssh -w 1 -c 3

```

## Setting

```
[plugin.checks.ssh-to-foo]
command = "/path/to/check-ssh -H foo.local -u foo -i /path/to/foo_id_rsa -t 10 -w 1 -c 10"
```
