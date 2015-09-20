docker-machine-xhyve
===

Docker Machine driver plugin for xhyve native OS X Hypervisor

## Install


Now, It is not yet available demonized.

It is as far to ssh login automatically.  
When docker-machine has finished creating a vm, at the same time xhyve also shut down.

```bash
# @nathanleclaire developpnig libmachine-rpc branch
go get github.com/nathanleclaire/machine
# Checkout branch
cd $GOPATH/src/github.com/nathanleclaire/machine
git checkout nathanleclaire/libmachine_rpc_plugins
# Make libmachine rpc include docker-machine_darwin-amd64 binary
script/build
# go get this repo
go get -d github.com/zchee/docker-machine-xhyve
# Intalll binary from /usr/local/bin/docker-machine-xhve
cd $GOPATH/src/github.com/zchee/docker-machine-xhyve
make install
# Create vm backend xhyve.
cd $GOPATH/src/github.com/nathanleclaire/machine
sudo ./docker-machine_darwin-amd64 create --driver xhyve xhyve
```

