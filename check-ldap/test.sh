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

# By default, LDAP_DOMAIN=example.org and LDAP_ORGANISATION=Example Inc
docker run --name test-$plugin -p 389:389 -d \
	-e 'LDAP_ADMIN_PASSWORD=passpass' \
	osixia/openldap
trap 'docker stop test-$plugin; docker rm test-$plugin; exit' EXIT
sleep 10

base_dn='dc=example,dc=org'
if $plugin -D "cn=admin,$base_dn" -P passpass -b "$base_dn" -c 2 -w 1
then
	echo OK
else
	echo FAIL
fi
