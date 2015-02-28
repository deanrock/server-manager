#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd $DIR

echo "DIR $DIR"

export LC_CTYPE=en_US.UTF-8
export LC_ALL=en_US.UTF-8

source env/bin/activate
python manage.py update-nginx-config --settings=manager.settings.$1
python manage.py runserver 0.0.0.0:5555 --settings=manager.settings.$1


