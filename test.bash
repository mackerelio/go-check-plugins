#!/bin/bash

# prepare tests
for f in check-*/test.sh check-*/test-*.sh
do
	dir=$(dirname "$f")
	name=$(basename "$dir")
	go build -o "$dir/$name" ./"$dir" || exit
done

# run tests
declare -A plugins=()
declare -a pids=()
for f in check-*/test.sh check-*/test-*.sh
do
	./"$f" &
	pid=$!
	plugins[$pid]="$f"
	pids+=("$pid")
done

# collect the results
declare -a results=()
status=0
for i in "${pids[@]}"
do
	if wait "$i"
	then
		results+=("OK: ${plugins[$i]}")
	else
		results+=("ERR: ${plugins[$i]}")
		status=1
	fi
done
echo '======' >&2
for s in "${results[@]}"
do
	echo "$s" >&2
done
exit $status
