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

docker run --name test-$plugin -p 16379:6379 -d redis:5 --requirepass passpass
trap 'docker stop test-$plugin; docker rm test-$plugin; exit' EXIT
sleep 10

if $plugin reachable --port 16379 --password passpass >/dev/null 2>&1
then
	echo OK
else
	echo FAIL
fi
