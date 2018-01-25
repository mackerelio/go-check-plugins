# check-uptime

## Description

Check uptime seconds.

## Setting

```
[plugin.checks.uptime]
command = "/path/to/check-uptime --warning-under=600 --critical-under=120"
```

## Options

```
  -w, --warning-under=N     Trigger a warning if under the seconds
  -c, --critical-under=N    Trigger a critial if under the seconds
  -W, --warning-over=N      Trigger a warning if over the seconds
  -C, --critical-over=N     Trigger a critical if over the seconds
  -r, --warning-rewind      Trigger a warning if rewind uptime (detect reboot)
  -R, --critical-rewind     Trigger a critical if rewind uptime (detect reboot)
  -s, --state-dir=DIR       Dir to keep state files under
```
