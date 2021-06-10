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

port=10389
image=osixia/openldap
password=passpass

# By default, LDAP_DOMAIN=example.org and LDAP_ORGANISATION=Example Inc
docker run --name "test-$plugin" -p "$port:389" -d \
	-e "LDAP_ADMIN_PASSWORD=$password" \
	"$image"
trap 'docker stop test-$plugin; docker rm test-$plugin; exit' EXIT
sleep 10

base_dn='dc=example,dc=org'
exec $plugin -p $port -D "cn=admin,$base_dn" -P "$password" -b "$base_dn" -c 2 -w 1
