#!/bin/sh

prog=$(basename "$0")
if ! [ -S /var/run/docker.sock ]
then
	echo "$prog: there are no running docker" >&2
	exit 2
fi

cd "$(dirname "$0")" || exit
plugin=$(basename "$(pwd)")
if ! which "$plugin" >/dev/null
then
	echo "$prog: $plugin is not installed" >&2
	exit 2
fi

password=passpass
port=13306
image=mysql:8
cacert=$(mktemp 'ca.XXXXXXXXXX.pem')

docker run -d \
	--name "test-$plugin" \
	-p "$port:3306" \
	-e "MYSQL_ROOT_PASSWORD=$password" "$image"
trap 'docker stop test-$plugin; docker rm test-$plugin; rm $cacert; exit 1' 1 2 3 15

# wait until bootstrap mysqld..
for i in $(seq 10)
do
	echo "Connecting $i..."
	if $plugin connection --port $port --password $password >/dev/null 2>&1
	then
		break
	fi
	sleep 3
done
sleep 1

docker cp "test-$plugin:/var/lib/mysql/ca.pem" "$cacert"

$plugin connection --port $port --password=$password \
	--tls --tls-root-cert="$cacert" --tls-skip-verify
status=$?
docker stop "test-$plugin"
docker rm "test-$plugin"
rm "$cacert"
exit $status
