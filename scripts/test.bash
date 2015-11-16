#!/bin/bash

export GOPATH=/Users/`users | awk 'NR == 1 {print $1}'`/go
export PATH=/usr/local/go/bin:$GOPATH/bin:$PATH
export GOOS=darwin
export GO15VENDOREXPERIMENT=1
export GOARCH=amd64
export GOMAXPROCS=8
export MACHINE_DEBUG=1
export MACHINE_DEBUG_DRIVER=1

make "$1"
