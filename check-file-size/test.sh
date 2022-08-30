#!/bin/sh

prog=$(basename "$0")
cd "$(dirname "$0")" || exit
PATH=$(pwd):$PATH
plugin=$(basename "$(pwd)")
if ! which "$plugin" >/dev/null
then
	echo "$prog: $plugin is not installed" >&2
	exit 2
fi

FILENAME=$(mktemp)
dd if=/dev/zero of=$FILENAME count=1023 bs=1
exec $plugin -b $FILENAME
