go-check-plugins
================

Check Plugins for monitoring written in golang.


Documentation for each plugin is located in its respective sub directory.

* [check-file-age](./check-file-age/README.md)
* [check-http](./check-http/README.md)
* [check-load](./check-load/README.md)
* [check-procs](./check-procs/README.md)

Specification
-------------

The specs for the check plugins are mostly the same as the [Nagios](https://www.nagios.org/) plugin. In the settings file, the assign command’s exit status will be treated as shown below.

| exit status          |  meaning |
|:---------------------|---------:|
| 0                    | OK       |
| 1                    | WARNING  |
| 2                    | CRITICAL |
| other than 0,1, or 2 | UNKNOWN  |


Contribution
------------

* fork it
* develop the plugin you want
* create a pullrequest!
