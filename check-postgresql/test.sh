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

user=postgres
password=passpass
port=15432
image=postgres:11

docker run -d \
	--name "test-$plugin" \
	-p "$port:5432" \
	-e "POSTGRES_PASSWORD=$password" "$image"
trap 'docker stop test-$plugin; docker rm test-$plugin; exit' EXIT
sleep 10

exec $plugin connection --port $port --user=$user --password=$password
