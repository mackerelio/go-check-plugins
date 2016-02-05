# check-jmx-jolokia

## Description

Checks for JMX value using jolokia.

## Setting

```
[plugin.checks.jmx_jolokia]
command = "/path/to/check-jmx-jolokia -H 127.0.0.1 -p 8778 -m java.lang:type=OperatingSystem -a ProcessCpuLoad -w 10 -c 20" 
```