#!/bin/bash

  sudo \
    env GOROOT="/Users/travis/.gimme/versions/go$TRAVIS_GO_VERSION.darwin.amd64" \
    env GOPATH="${TRAVIS_BUILD_DIR}/Godeps/_workspace:$GOPATH"
    env PATH=$TRAVIS_BUILD_DIR/Godeps/_workspace/bin:/Users/travis/gopath/bin:$GOROOT/bin:/Users/travis/bin:/Users/travis/bin:$PATH \
    env GO15VENDOREXPERIMENT=1 \
    env GOMAXPROCS=8 \
    env MACHINE_DEBUG=1 \
  "$@"
