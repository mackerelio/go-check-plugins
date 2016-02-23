# check-file-size

## Description

Check file size in specified directory.

## Setting

```
[plugin.checks.filesize]
command = "/path/to/check-file-size -b /fluentd/buffer_dir/ -w 5M -c 10M"
```

## Options

```
-b --base     The base directory(required)
-w --warning  The warning threshold of file size (default: 1K)
-c --critical The critical threshold of file size (default: 1K)
-d --depth    Max depth of the directory from base directory (default: 1)
```
