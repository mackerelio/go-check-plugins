#!/bin/sh

prog=$(basename "$0")
if ! [ -S /var/run/docker.sock ]
then
	echo "$prog: there are no running docker" >&2
	exit 2
fi

cd "$(dirname "$0")" || exit
PATH=$(pwd):$PATH
plugin=$(basename "$(pwd)")
if ! which "$plugin" >/dev/null
then
	echo "$prog: $plugin is not installed" >&2
	exit 2
fi

port=21211
image=memcached

docker run --name "test-$plugin" -p "$port:11211" -d "$image"
trap 'docker stop test-$plugin; docker rm test-$plugin; exit' EXIT
sleep 10

exec $plugin -p $port -k test
