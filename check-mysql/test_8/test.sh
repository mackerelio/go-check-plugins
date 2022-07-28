#!/bin/sh

set -ex

echo $0
prog=$(basename $0)
if ! command -v docker-compose > /dev/null
then
	echo "$prog: docker-compose is not installed" >&2
	exit 2
fi

cd $(dirname $0)
plugin=$(basename $(realpath ../))
if ! which "$plugin" >/dev/null
then
	echo "$prog: $plugin is not installed" >&2
	exit 2
fi

user=root
password=mysql
primary_port=13306
replica_port=23306

if $plugin connection --host=127.0.0.1 --port=$primary_port --user=$user --password=$password; then
	echo 'FAIL: connection shoule not be OK'
	exit 1
fi

docker-compose up -d
trap 'docker-compose down --rmi local -v' EXIT

USER=$user PASSWORD=$password PORT=$primary_port ../wait.sh

if ! $plugin connection --host=127.0.0.1 --port=$primary_port --user=$user --password=$password; then
	echo 'FAIL: connection should be OK'
	exit 1
fi

sleep 2

if ! $plugin uptime --host=127.0.0.1 --port=$primary_port --user=$user --password=$password --critical=2 --warning=1; then
	echo 'FAIL: uptime should be OK'
	exit 1
fi

if $plugin uptime --host=127.0.0.1 --port=$primary_port --user=$user --password=$password --critical=1000000000; then
	echo 'FAIL: uptime should not be OK'
	exit 1
fi

if ! $plugin readonly --host=127.0.0.1 --port=$primary_port --user=$user --password=$password OFF; then
	echo 'FAIL: primary server should not be read-only'
	exit 1
fi

if ! $plugin readonly --host=127.0.0.1 --port=$replica_port --user=$user --password=$password ON; then
	echo 'FAIL: replica server should be read-only'
	exit 1
fi

mysql -u$user -p$password --host 127.0.0.1 --port=$replica_port -e """
CHANGE REPLICATION SOURCE TO SOURCE_HOST='primary', SOURCE_PORT=3306;
"""

if ! $plugin replication --host=127.0.0.1 --port=$primary_port --user=$user --password=$password; then
	# MySQL is not a replica (primary server)
	echo 'FAIL: replication of primary server should be OK'
	exit 1
fi

if $plugin replication --host=127.0.0.1 --port=$replica_port --user=$user --password=$password; then
	# Replication does not start yet
	echo 'FAIL: replication of replica server should not be started'
	exit 1
fi

mysql -u$user -p$password --host 127.0.0.1 --port=$replica_port -e """
START REPLICA USER='repl' PASSWORD='repl';
"""

# starting replication may take a while
sleep 1

if ! $plugin replication --host=127.0.0.1 --port=$replica_port --user=$user --password=$password; then
	echo 'FAIL: replication of replica server should be started'
	exit 1
fi

echo OK
