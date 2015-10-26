go-check-plugins
================

Check Plugins for monitoring written in golang.


Documentation for each plugin is located in its respective sub directory.

* [check-file-age](./check-file-age/README.md)
* [check-http](./check-http/README.md)
* [check-load](./check-load/README.md)
* [check-mysql](./check-mysql/README.md)
* [check-procs](./check-procs/README.md)

Specification
-------------

The specs for the check plugins are mostly the same as the plugins for [Nagios](https://www.nagios.org/) and [Sensu](https://sensuapp.org/).
The exit status of the commands are treated as follows.

| exit status           |  meaning |
|:----------------------|---------:|
| 0                     | OK       |
| 1                     | WARNING  |
| 2                     | CRITICAL |
| other than 0, 1, or 2 | UNKNOWN  |


Installation
------------

Install the plugin package from either the yum or the apt repository.

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
* create a pull request!
