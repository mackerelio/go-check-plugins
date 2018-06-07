# check-mailq

## Description

Monitor mail queue count.

## Usage

```
$ check-mailq -w 100 -c 200 -M postfix
```

## Setting

```
[plugin.checks.mailq]
command = "/path/to/check-mailq -w 100 -c 200 -M postfix"
```

## Other

* inspired by [check_mailq.pl](https://github.com/nagios-plugins/nagios-plugins/blob/master/plugins-scripts/check_mailq.pl)
