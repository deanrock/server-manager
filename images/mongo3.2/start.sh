#!/bin/bash

mkdir -p /home/#user#/apps/#appname#/data/db/
mkdir -p /home/#user#/apps/#appname#/log/

FILE=/home/#user#/apps/#appname#/.mongodb_admin_password
PASS=$(pwgen -s 20 1)

CMD="mongod -f /etc/mongod.conf --storageEngine wiredTiger --master --auth"

if [ ! -f $FILE ]; then
	$CMD &

	RET=1
	while [[ RET -ne 0 ]]; do
	    echo "=> Waiting for confirmation of MongoDB service startup"
	    sleep 5
	    mongo admin --eval "help" >/dev/null 2>&1
	    RET=$?
	done

	echo "=> Creating an admin user in MongoDB"
	mongo admin --eval "db.createUser({user: 'admin', pwd: '$PASS', roles:[{role:'root',db:'admin'}]});"

	echo "=> Password for admin user written to $FILE!"
	echo $PASS > $FILE

	echo "=> Sleeping for 5 seconds"
	sleep 5

	echo "=> Killing mongod instance"
	killall -9 mongod
fi

$CMD
