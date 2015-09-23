go-check-plugins
================

Check Plugins for monitoring written in golang.

Specification
-------------

The specs for the check plugins are mostly the same as the [Nagios](https://www.nagios.org/) plugin. In the settings file, the assign command’s exit status will be treated as shown below.

| exit status          |  meaning |
|:---------------------|---------:|
| 0                    | OK       |
| 1                    | WARNING  |
| 2                    | CRITICAL |
| other than 0,1, or 2 | UNKNOWN  |
