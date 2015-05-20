# check-file-age

## Description

Monitor file age and size for script monitoring.

## Usage

```
$ check-file-age -w 240 -W 10 -c 600 -C 0 -f filename
```

## Setting

```
[plugin.checks.fileage]
command = "/path/to/check-file-age -f filename"
```

## Other

* inspired by [check_file_age.pl](https://github.com/nagios-plugins/nagios-plugins/blob/master/plugins-scripts/check_file_age.pl)
