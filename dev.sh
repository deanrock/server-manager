#!/bin/bash

#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd $DIR

export GOROOT=/usr/local/go1.4/go
export GOPATH=$DIR/go-libs/
export GOBIN=$DIR/go-libs/bin/

export LC_CTYPE=en_US.UTF-8
export LC_ALL=en_US.UTF-8

if [ $1 = "cron" ]; then
	echo "cron"

	cd ./cron
	/usr/local/go1.4/go/bin/go run main.go
elif [ $1 = "uptime-updater" ]; then
	echo "uptime updater"

	cd ./uptime-updater
	/usr/local/go1.4/go/bin/go run main.go
elif [ $1 = "proxy" ]; then
	echo "proxy"

	cd ./proxy
	/usr/local/go1.4/go/bin/go run main.go session.go
elif [ $1 = "ssh" ]; then
	echo "ssh"

	cd ./ssh-server
	/usr/local/go1.4/go/bin/go run main.go
elif [ $1 = "ondemand" ]; then
	echo "ondemand"

	cd ./ondemand
	/usr/local/go1.4/go/bin/go run main.go
elif [ $1 = "manager" ]; then
	echo "manager"

	source env/bin/activate
	python manage.py runserver 0.0.0.0:5555 --settings=manager.settings.dev --noreload
else
	echo "wrong argument"
fi
