#!/bin/bash
#!/bin/bash

green='\033[0;32m'
NC='\033[0m' # No Color

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

sudo rm /usr/bin/shell #remove old python shell if still exists
sudo rm /usr/bin/manager-shell # remove old golang shell if still exists


#golang build stuff

cd $DIR

mkdir -p ./bin

export GOROOT=/usr/local/go1.4/go
export GOPATH=$DIR/go-libs/
export GOBIN=$DIR/go-libs/bin/

echo -e "${green}[go] compiling proxy ...${NC}"
cd $DIR/proxy

echo -e "${green}[go] getting libraries ...${NC}"
/usr/local/go1.4/go/bin/go get
/usr/local/go1.4/go/bin/go install
/usr/local/go1.4/go/bin/go build -o ../bin/proxy main.go session.go

echo -e "${green}[go] compiling cron ...${NC}"
cd $DIR/cron

echo -e "${green}[go] getting libraries ...${NC}"
/usr/local/go1.4/go/bin/go get
/usr/local/go1.4/go/bin/go install
/usr/local/go1.4/go/bin/go build -o ../bin/cron main.go

echo -e "${green}[go] compiling ssh server ...${NC}"
cd $DIR/ssh-server

echo -e "${green}[go] getting libraries ...${NC}"
/usr/local/go1.4/go/bin/go get
/usr/local/go1.4/go/bin/go install
/usr/local/go1.4/go/bin/go build -o ../bin/ssh-server main.go

echo -e "${green}[all] finished${NC}"
