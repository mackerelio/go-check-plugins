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

# GitHub-hosted Linux runner mounts tracefs on /sys/kernel/debug/tracing.
# That fstype don't appear in ME_DUMMY macro on coreutils.
# Thus, for now, we ignore it only in test.sh.
exec $plugin -X tracefs
