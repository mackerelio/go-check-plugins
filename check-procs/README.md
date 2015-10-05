# check-procs
## Description

Checks processes if the specific metric is outside the required threshold ranges.

## Setting

```
[plugin.checks.procs]
command = "/path/to/check-procs --pattern=PROCESS_NAME --state=STATE --warn-under=N"
```

## Other

* This is a Go port of [Sensu-Plugins-process-checks](https://github.com/sensu-plugins/sensu-plugins-process-checks).
* [Nagios Plugins - check_procs.c](https://github.com/nagios-plugins/nagios-plugins/blob/master/plugins/check_procs.c)
