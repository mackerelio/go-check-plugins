# check-disk

## Description

Check free space of disk.

## Setting

```
[plugin.checks.disk]
command = "/path/to/check-disk --warning=10 --critical=5 --path=/"
```

## Options

```
-w, --warning=N, N%     Exit with WARNING status if less than N units or N% of disk are free
-c, --critical=N, N%    Exit with CRITICAL status if less than N units or N% of disk are free
-p, --path=PATH         Mount point or block device as emitted by the mount(8) command
-u, --units=STRING      Choose bytes, kB, MB, GB, TB (default: MB)
```
