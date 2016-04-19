# check-procs
## Description

Checks processes if the specific metric is outside the required threshold ranges.

## Setting

```
[plugin.checks.procs]
command = "/path/to/check-procs --pattern=PROCESS_NAME --state=STATE --warning-under=N"
```
## Options

```
-w, --warning-over=N                Trigger a warning if over a number
-c, --critical-over=N               Trigger a critical if over a number
-W, --warning-under=N               Trigger a warning if under a number (default: 1)
-C, --critical-under=N              Trigger a critial if under a number (default: 1)
-m, --match-self                    Match itself
-M, --match-parent                  Match parent
-p, --pattern=PATTERN               Match a command against this pattern
-x, --exclude-pattern=PATTERN       Don't match against a pattern to prevent false positives
    --ppid=PPID                     Check against a specific PPID
-f, --file-pid=PID                  Check against a specific PID
-z, --virtual-memory-size=VSZ       Trigger on a Virtual Memory size is bigger than this
-r, --resident-set-size=RSS         Trigger on a Resident Set size is bigger than this
-P, --proportional-set-size=PCPU    Trigger on a Proportional Set Size is bigger than this
-T, --thread-count=THCOUNT          Trigger on a Thread Count is bigger than this
-s, --state=STATE                   Trigger on a specific state, example: Z for zombie
-u, --user=USER                     Trigger on a specific user
-U, --user-not=USER                 Trigger if not owned a specific user
-e, --esec-over=SECONDS             Match processes that older that this, in SECONDS
-E, --esec-under=SECONDS            Match process that are younger than this, in SECONDS
-i, --cpu-over=SECONDS              Match processes cpu time that is older than this, in SECONDS
-I, --cpu-under=SECONDS             Match processes cpu time that is younger than this, in SECONDS
```

## Other

* This is a Go port of [Sensu-Plugins-process-checks](https://github.com/sensu-plugins/sensu-plugins-process-checks).
* [Nagios Plugins - check_procs.c](https://github.com/nagios-plugins/nagios-plugins/blob/master/plugins/check_procs.c)
