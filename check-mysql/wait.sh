#!/bin/sh

# wait until bootstrap mysqld.. (max = 60s)
for i in $(seq 20)
do
	echo "Connecting $i..."
	if mysql -u"$USER" -p"$PASSWORD" --host=127.0.0.1 --port="$PORT" -e"select 1" >/dev/null; then
		break
	fi
	sleep 3
done
sleep 1
