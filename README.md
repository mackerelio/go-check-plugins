go-check-plugins
================

Check Plugins for monitoring written in golang.


Documentation for each plugin is located in its respective sub directory.

* [check-aws-sqs-queue-size](./check-aws-sqs-queue-size/README.md)
* [check-cert-file](./check-cert-file/README.md)
* [check-disk](./check-disk/README.md)
* [check-elasticsearch](./check-elasticsearch/README.md)
* [check-file-age](./check-file-age/README.md)
* [check-file-size](./check-file-size/README.md)
* [check-http](./check-http/README.md)
* [check-jmx-jolokia](./check-jmx-jolokia/README.md)
* [check-load](./check-load/README.md)
* [check-log](./check-log/README.md)
* [check-mailq](./check-mailq/README.md)
* [check-masterha](./check-masterha/README.md)
* [check-memcached](./check-memcached/README.md)
* [check-mysql](./check-mysql/README.md)
* [check-ntpoffset](./check-ntpoffset/README.md)
* [check-ntservice](./check-ntservice/README.md)
* [check-postgresql](./check-postgresql/README.md)
* [check-procs](./check-procs/README.md)
* [check-redis](./check-redis/README.md)
* [check-solr](./check-solr/README.md)
* [check-ssh](./check-ssh/README.md)
* [check-tcp](./check-tcp/README.md)
* [check-uptime](./check-uptime/README.md)
* [check-windows-eventlog](./check-windows-eventlog/README.md)

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
To setup these package repositories, see the documentation regarding the installation of mackerel-agent ([rpm](https://mackerel.io/docs/entry/howto/install-agent/rpm) / [deb](https://mackerel.io/docs/entry/howto/install-agent/deb)).

mackerel-check-plugins will be installed to ```/usr/local/bin/check-*```.

### yum

```shell
yum install mackerel-check-plugins
```

### apt

```shell
apt-get install mackerel-check-plugins
```

Use check plugins in Mackerel
-----------------------------

See the following documentation.

English: https://mackerel.io/docs/entry/custom-checks

Japanese: https://mackerel.io/ja/docs/entry/custom-checks


Contribution
------------

* fork it
* develop the plugin you want
* create a pull request!
