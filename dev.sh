#!/bin/bash

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd $DIR
export GOROOT=/usr/local/go1.4/go
export GOPATH=$DIR/go-libs/
export GOBIN=$DIR/go-libs/bin/

if [ $1 = "ssh" ]; then
	if [ ! -f ./id_rsa ]; then
		ssh-keygen -b 2048 -t rsa -f ./id_rsa -q -N ""
	fi
fi

/usr/local/go1.4/go/bin/go run main.go $1
