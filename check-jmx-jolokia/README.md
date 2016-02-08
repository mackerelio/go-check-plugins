# check-jmx-jolokia

## Description

Checks for JMX value using jolokia.

## Setting

```
[plugin.checks.jmx_jolokia]
command = "/path/to/check-jmx-jolokia -H 127.0.0.1 -p 8778 -m java.lang:type=OperatingSystem -a ProcessCpuLoad -w 10 -c 20"
```

## Options
```
-H, --host=       Host name or IP Address
-p, --port=       Port (default: 8778)
-t, --timeout=    Seconds before connection times out (default: 10)
-m, --mbean=      MBean
-a, --attribute=  Attribute
-i, --inner-path= InnerPath
-k, --key=        Key (default: value)
-w, --warning=    Trigger a warning if over a number
-c, --critical=   Trigger a critical if over a number
```