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

http_port=10080
http_image=httpd:2.4-bullseye
proxy_port=13128
proxy_image=ubuntu/squid:5.2-22.04_beta

halt()
{
	docker stop "test-$plugin" "test-$plugin-proxy"
	docker rm "test-$plugin" "test-$plugin-proxy"
	docker network rm "test-$plugin-net"
}

docker network create --driver bridge "test-$plugin-net"
docker run --name "test-$plugin" \
	--net "test-$plugin-net" -p "$http_port:80" \
	-v "$(pwd)/testdata:/usr/local/apache2/htdocs/" -d "$http_image"
docker run --name "test-$plugin-proxy" \
	--net "test-$plugin-net" -p "$proxy_port:3128" \
	-d "$proxy_image"
trap 'halt; exit 1' 1 2 3 15
sleep 10

status=0
$plugin -u "http://localhost:$http_port/"
((status+=$?))

# Target hostname in the URL will be resolved in the side of proxy container,
# thus it must be container name or its ID.
$plugin -u "http://test-$plugin/" --proxy="http://localhost:$proxy_port"
((status+=$?))

halt
exit "$status"
