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
-w, --warning=N, N%                  Exit with WARNING status if less than N units or N% of disk are free
-c, --critical=N, N%                 Exit with CRITICAL status if less than N units or N% of disk are free
-W, --iwarning=N%                    Exit with WARNING status if less than PERCENT of inode space is free
-K, --icritical=N%                   Exit with CRITICAL status if less than PERCENT of inode space is free
-p, --path=PATH                      Mount point or block device as emitted by the mount(8) command (may be repeated)
-x, --exclude-device=EXCLUDE PATH    Ignore device (only works if -p unspecified)
-A, --all                            Explicitly select all paths.
-X, --exclude-type=TYPE              Ignore all filesystems of indicated type (may be repeated)
-N, --include-type=TYPE              Check only filesystems of indicated type (may be repeated)
-u, --units=STRING                   Choose bytes, kB, MB, GB, TB (default: MB)
```
