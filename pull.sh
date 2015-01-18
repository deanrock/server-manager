#!/bin/bash
pip install -r requirements.txt

cp manager/db.sqlite3 ./backup_db_`date +"%Y-%m-%d_%H-%M-%S"`.db
python manage.py migrate --settings=$1
python manage.py syncimageconfig --settings=$1

sudo cp shell/shell.py /usr/bin/shell
