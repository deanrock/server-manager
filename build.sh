#!/bin/bash
set -e

green='\033[0;32m'
NC='\033[0m' # No Color

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd $DIR

echo "DIR $DIR"

export LC_CTYPE=en_US.UTF-8
export LC_ALL=en_US.UTF-8

#golang build stuff

cd $DIR

mkdir -p ./bin

export GOPATH=$DIR/go-libs/
export GOBIN=$DIR/go-libs/bin/

echo -e "${green}[go] getting libraries ...${NC}"
go get
go install
go build -o ./bin/server-manager main.go

echo -e "${green}[go] archiving files ...${NC}"
tar cvfz package.tar.gz bin/ config-example.json static/ proxy/templates/ images/

echo -e "${green}[all] finished${NC}"
