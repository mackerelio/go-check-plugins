# check-disk

## Description

Check free space of disk.

## Setting

```
[plugin.checks.disk]
command = "/path/to/check-disk --warning=20 --critical=10 --path=/"
```

## Options

```
-w, --warning=N     Exit with WARNING status if less than N GB of disk are free
-c, --critical=N    Exit with CRITICAL status if less than N GB of disk are free
-p, --path=PATH     Mount point or block device as emitted by the mount(8) command
```
