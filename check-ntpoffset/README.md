# check-ntpoffset

## Description

Check ntp offset.

## Setting

```
[plugin.checks.ntpoffset]
command = "/path/to/check-ntpoffset -w=50 -c=100"
```

## Options

```
-w --warning   Warning threshold of ntp offset(ms) (default: 50)
-c --critical  Critical threshold of ntp offset(ms) (default: 100)
```


