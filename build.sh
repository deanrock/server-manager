#!/bin/bash
set -e

green='\033[0;32m'
NC='\033[0m' # No Color

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd $DIR

echo "DIR $DIR"

#golang build stuff
mkdir -p ./bin

export GOPATH=$DIR/go-libs/
export GOBIN=$DIR/go-libs/bin/

echo -e "${green}[go] getting libraries ...${NC}"
go build -o ./bin/server-manager main.go

echo -e "${green}[go] archiving files ...${NC}"
tar cvfz package.tar.gz bin/ config-example.json static/ proxy/templates/ images/

echo -e "${green}[all] finished${NC}"
