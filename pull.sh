#!/bin/bash
#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd $DIR

echo "DIR $DIR"

export LC_CTYPE=en_US.UTF-8
export LC_ALL=en_US.UTF-8

if [ -d "./env" ]; then
    echo "env already exists"
else
   virtualenv ./env
fi

source env/bin/activate
pip install -r requirements.txt

cp manager/db.sqlite3 ./backup_db_`date +"%Y-%m-%d_%H-%M-%S"`.db
python manage.py migrate --settings=manager.settings.$1
python manage.py syncimageconfig --settings=manager.settings.$1

sudo cp shell/shell.py /usr/bin/shell



#build proxy

cd $DIR

mkdir -p ./bin

cd ./proxy

export GOROOT=/usr/local/go1.4/go
export GOPATH=$DIR/go-libs/
export GOBIN=$DIR/go-libs/bin/

/usr/local/go1.4/go/bin/go get
/usr/local/go1.4/go/bin/go install
/usr/local/go1.4/go/bin/go build -o ../bin/proxy main.go session.go shell.go
