#!/bin/sh

prog=$(basename $0)
cd $(dirname $0)
plugin=$(basename $(pwd))
if ! which -s $plugin
then
	echo "$prog: $plugin is not installed" >&2
	exit 2
fi

if $plugin >/dev/null 2>&1; then
	echo OK
else
	echo FAIL
fi
