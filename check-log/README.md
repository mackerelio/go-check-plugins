# check-log

## Description

Checks a log file using a regular expression.

## Setting

```
[plugin.checks.log]
command = "/path/to/check-log --file=/path/to/file --pattern=REGEXP --warning-over=N --critical-over=N"
```

To specify encoding of the log files, you can use `--encoding` option. Below's list of supported encodings.

* UTF-8
* CP437
* CP866
* ISO-2022-JP
* LATIN-1
* ISO-8859-1
* ISO-8859-2
* ISO-8859-3
* ISO-8859-4
* ISO-8859-5
* ISO-8859-6
* ISO-8859-7
* ISO-8859-8
* ISO-8859-10
* ISO-8859-13
* ISO-8859-14
* ISO-8859-15
* ISO-8859-16
* KOI8R
* KOI8U
* Macintosh
* MacintoshCyrillic
* Windows1250
* Windows1251
* Windows1252
* Windows1253
* Windows1254
* Windows1255
* Windows1256
* Windows1257
* Windows1258
* Windows874
* XUserDefined
* Big5
* EUC-KR
* HZ-GB2312
* sjis
* CP932
* Shift_JIS
* EUC-JP
* UTF-16 (detect BOM)
* UTF-16BE
* UTF-16LE

## See Other

* inspired by [sensu-plugins-logs](https://github.com/sensu-plugins/sensu-plugins-logs).
