# check-procs

## Description
Checks processes if the specific metric is outside the required threshold ranges.

## Synopsis
```
check-procs --pattern=PROCESS_NAME --state=STATE --warning-under=N
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-procs
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-procs --pattern=PROCESS_NAME --state=STATE --warning-under=N
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.check-procs-sample]
command = ["check-procs", "--pattern", "<PROCESS_NAME>", "--state", "<STATE>", "--warning-under", "N"]
```

## Usage
### Options

```
  -w, --warning-over=N                Trigger a warning if over a number
      --warn-over=N                   (DEPRECATED) Trigger a warning if over a number
  -c, --critical-over=N               Trigger a critical if over a number
  -W, --warning-under=N               Trigger a warning if under a number (default: 1)
      --warn-under=N                  (DEPRECATED) Trigger a warning if under a number (default: 1)
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

## For more information
Please refer to the following.

- Read Mackerel Docs; [Monitoring Processes - Mackerel Docs](https://mackerel.io/docs/entry/howto/check/process)
- Execute `check-procs -h` and you can get command line options.

## Other

- This is a Go port of [Sensu-Plugins-process-checks](https://github.com/sensu-plugins/sensu-plugins-process-checks).
- [Nagios Plugins - check_procs.c](https://github.com/nagios-plugins/nagios-plugins/blob/master/plugins/check_procs.c)
