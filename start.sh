#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd $DIR

echo "DIR $DIR"

source env/bin/activate
python manage.py update-nginx-config --settings=manager.settings.$1
python manage.py runserver 0.0.0.0:4444 --settings=manager.settings.$1


