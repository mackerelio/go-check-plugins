go-check-plugins
================

Check Plugins for monitoring written in golang.


Documentation for each plugin is located in its respective sub directory.

* [check-file-age](./check-file-age/README.md)
* [check-http](./check-http/README.md)
* [check-load](./check-load/README.md) (This plugin is not yet included within the yum/apt packages)
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


Installation
------------

Install the plugin pack from either the yum or the apt repository.

### CentOS 5/6

```shell
yum install mackerel-check-plugins
```

### Debian 6/7

```shell
apt-get install mackerel-check-plugins
```

mackerel-check-plugins will be installed to ```/usr/local/bin/check-*```.


Use check plugins in Mackerel
-----------------------------

See the following documentation.

English: http://help.mackerel.io/entry/custom-checks

Japanese: http://help-ja.mackerel.io/entry/custom-checks


Contribution
------------

* fork it
* develop the plugin you want
* create a pullrequest!
