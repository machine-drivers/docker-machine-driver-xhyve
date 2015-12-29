#!/bin/bash

sudo \
  env GOROOT=$GOROOT \
  env GOPATH="$GOPATH" \
  env PATH=$PATH \
  env GO15VENDOREXPERIMENT=1 \
  env GOMAXPROCS=8 \
  env MACHINE_DEBUG=1 \
"$@"
