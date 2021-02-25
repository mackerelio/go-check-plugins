#!/bin/sh

prog=$(basename $0)
if ! [[ -S /var/run/docker.sock ]]
then
	echo "$prog: there are no running docker" >&2
	exit 2
fi

cd $(dirname $0)
plugin=$(basename $(pwd))
if ! which -s $plugin
then
	echo "$prog: $plugin is not installed" >&2
	exit 2
fi

sleep_time=20
user=root
password=mysql
port=13306

if $plugin connection --host=127.0.0.1 --port=$port --user=$user --password=$password; then
	echo 'FAIL: connection result shoule not be OK'
	exit 1
fi

docker run -d \
	--name test-$plugin \
	-p $port:3306 \
	-e MYSQL_ROOT_PASSWORD=$password \
	mysql:8
trap 'docker stop test-$plugin; docker rm test-$plugin; exit' EXIT
sleep $sleep_time

if ! $plugin connection --host=127.0.0.1 --port=$port --user=$user --password=$password; then
	echo 'FAIL: connection should be OK'
	exit 1
fi

if ! $plugin uptime --host=127.0.0.1 --port=$port --user=$user --password=$password --critical=2 --warning=1; then
	echo 'FAIL: uptime should be OK'
	exit 1
fi

if $plugin uptime --host=127.0.0.1 --port=$port --user=$user --password=$password --critical=1000000000; then
	echo 'FAIL: uptime should not be OK'
	exit 1
fi

if ! $plugin readonly --host=127.0.0.1 --port=$port --user=$user --password=$password OFF; then
	echo 'FAIL: readonly should be OK'
	exit 1
fi

if $plugin readonly --host=127.0.0.1 --port=$port --user=$user --password=$password ON; then
	echo 'FAIL: readonly should not be OK'
	exit 1
fi

if ! $plugin replication --host=127.0.0.1 --port=$port --user=$user --password=$password ON; then
	echo 'FAIL: replication should be OK'
	exit 1
fi

echo OK
