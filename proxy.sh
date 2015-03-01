#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd $DIR
cd ./proxy

export GOROOT=/usr/local/go1.4/go
export GOPATH=$DIR/go-libs/
export GOBIN=$DIR/go-libs/bin/

/usr/local/go1.4/go/bin/go get
/usr/local/go1.4/go/bin/go install
/usr/local/go1.4/go/bin/go run main.go session.go shell.go

