#!/bin/bash
#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd $DIR

echo "DIR $DIR"

if [ -d "./env" ]; then
    echo "env already exists"
else
   virtualenv ./env
fi

source env/bin/activate
pip install -r requirements.txt

cp manager/db.sqlite3 ./backup_db_`date +"%Y-%m-%d_%H-%M-%S"`.db
python manage.py migrate --settings=$1
python manage.py syncimageconfig --settings=$1

sudo cp shell/shell.py /usr/bin/shell
