# check-log

## Description

Checks a log file using a regular expression.

## Synopsis
```
check-log --file=/path/to/file --pattern=REGEXP --warning-over=N --critical-over=N
```

## Installation

First, build this program.

```
go get github.com/mackerelio/go-check-plugins
cd $(go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-log
go install
```

Or you can use this program by installing the official Mackerel package. See [Using the official check plugin pack for check monitoring - Mackerel Docs](https://mackerel.io/docs/entry/howto/mackerel-check-plugins).


Next, you can execute this program :-)

```
check-log --file=/path/to/file --pattern=REGEXP --warning-over=N --critical-over=N
```


## Setting for mackerel-agent

If there are no problems in the execution result, add a setting in mackerel-agent.conf .

```
[plugin.checks.check-log-sample]
command = ["check-log", "--file", "/path/to/file", "--pattern", "REGEXP", "--warning-over", "N", "--critical-over", "N"]
```

## Usage
### Options

```
  -f, --file=FILE                                Path to log file
  -p, --pattern=PAT                              Pattern to search for. If specified multiple, they will be treated together with the AND operator
      --suppress-pattern                         Suppress pattern display
  -E, --exclude=PAT                              Pattern to exclude from matching
  -w, --warning-over=                            Trigger a warning if matched lines is over a number
  -c, --critical-over=                           Trigger a critical if matched lines is over a number
      --warning-level=N                          Warning level if pattern has a group
      --critical-level=N                         Critical level if pattern has a group
  -r, --return                                   Return matched line
  -F, --file-pattern=FILE                        Check a pattern of files, instead of one file
  -i, --icase                                    Run a case insensitive match
  -s, --state-dir=DIR                            Dir to keep state files under
      --no-state                                 Don't use state file and read whole logs
      --encoding=                                Encoding of log file
      --missing=(CRITICAL|WARNING|OK|UNKNOWN)    Exit status when log files missing (default: UNKNOWN)
      --check-first                              Check the log on the first run
```

#### Using glob

You can check multiple files by using globs (and zsh extented globs by [mattn/go-zglob](https://github.com/mattn/go-zglob)) in `--file` option.
For example, `--file=/tmp/some.log_*` will check all of `/tmp/some.log_1`, `/tmp/some.log_2`, and so on.

And since `command` string in mackerel-agent.conf will be parsed by shell (in *nix `/bin/sh -c`, in Windows `cmd /c`), specifying glob like `--file /tmp/some.log_*` does not work as expected.
It will be expanded like `--file /tmp/some.log_1 /tmp/some.log_2`, so it will check only `/tmp/some.log_1`.

Therefore, when you want to check multiple files, use `--file=<glob>`, not `--file <glob>`, or please specify `command` by array.

#### Encoding

To specify encoding of the log files, you can use `--encoding` option. Below's list of supported encodings.

- UTF-8
- CP437
- CP866
- ISO-2022-JP
- LATIN-1
- ISO-8859-1
- ISO-8859-2
- ISO-8859-3
- ISO-8859-4
- ISO-8859-5
- ISO-8859-6
- ISO-8859-7
- ISO-8859-8
- ISO-8859-10
- ISO-8859-13
- ISO-8859-14
- ISO-8859-15
- ISO-8859-16
- KOI8R
- KOI8U
- Macintosh
- MacintoshCyrillic
- Windows1250
- Windows1251
- Windows1252
- Windows1253
- Windows1254
- Windows1255
- Windows1256
- Windows1257
- Windows1258
- Windows874
- XUserDefined
- Big5
- EUC-KR
- HZ-GB2312
- sjis
- CP932
- Shift_JIS
- EUC-JP
- UTF-16 (detect BOM)
- UTF-16BE
- UTF-16LE

## For more information
Please refer to the following.

- Read Mackerel Docs; [Monitoring Logs - Mackerel Docs](https://mackerel.io/docs/entry/howto/check/log)
- Execute `check-log -h` and you can get command line options.

## other
- inspired by [sensu-plugins-logs](https://github.com/sensu-plugins/sensu-plugins-logs).
