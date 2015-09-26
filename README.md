docker-machine-xhyve
===

Docker Machine driver plugin for xhyve native OS X Hypervisor

Master branch inherited from [nathanleclaire/docker-machine-xhyve](https://github.com/nathanleclaire/docker-machine-xhyve). Thanks [@nathanleclaire](https://github.com/nathanleclaire) :)  
If you have issues or pull-requests, Desired to be posted to this repository.


## Required

### docker-machine
**!! Please do not post the issue of this repository to the docker/machine !!**  
It will interfere with the development of the docker-machine.  
If you were doubt problem either, please post to this repository. I will judge.

Now, docker-machine are develop to new driver plugin mechanism.  
`docker-machine-xhyve` is using it.  
So, please try [nathanleclaire/machine/libmachine_rpc_plugins](https://github.com/nathanleclaire/machine/tree/libmachine_rpc_plugins) branch.

```bash
# @nathanleclaire developpnig libmachine-rpc branch
go get github.com/nathanleclaire/machine
cd $GOPATH/src/github.com/nathanleclaire/machine
git checkout nathanleclaire/libmachine_rpc_plugins
# Make libmachine rpc include docker-machine_darwin-amd64 binary
script/build
```

### xhyve-bindings
Since it is was hard to `os.exec` itself that embedded `xhyve.Exec`, for the time being, is separated into [xhyve-bindings](https://github.com/zchee/xhyve-bindings/tree/daemonize).

```bash
$ go get -d github.com/zchee/xhyve-bindings
$ cd $GOPATH/src/github.com/zchee/xhyve-bindings
$ git checkout daemonize
$ make
$ make install
```

#### boot2docker custom ISO
For now, using custom boot2docker ISO.  
If you want to know custom point, See https://github.com/zchee/boot2docker-legacy/tree/xhyve

```bash
$ git clone https://github.com/zchee/boot2docker-legacy.git
$ cd boot2docker-legacy
$ git checkout xhyve
$ docker build -t boot2docker-xhyve .
$ docker run --rm boot2docker-xhyve > boot2docker.iso
$ mv ./boot2docker.iso ~/.docker/machine/cache/
```


## Install

```bash
$ go get -d github.com/zchee/docker-machine-xhyve
$ cd $GOPATH/github.com/zchee/docker-machine-xhyve
$ make
$ make install
```


## TODO

- [x] Daemonize xhyve use `syscall` or `go execute external process myself` or `OS X launchd daemon` or other daemonize method
    - Since it is was hard to exec itself that embedded xhyve.Exec, for the time being, is separated into [xhyve-bindings](https://github.com/zchee/xhyve-bindings/tree/daemonize).

- [ ] Shared folder support

- [ ] Replace execute binary to syscall of golang
    - e.g. `hdutil`, ~~`dd`~~
    - Create blank disk use `libguestfs` instead of `dd`

- [ ] Replace generate uuid, execute `uuidgen` to native golang

- [x] Replace exec uuid2mac binary to standalone `vmnet.go`, `dhcp.go`

- [x] ~~Update xhyve source to unofficial edge branch~~
    - ~~See [update-xhyve-to-edge](https://github.com/zchee/docker-machine-xhyve/tree/update-xhyve-to-edge)~~
    - ~~Replace `Grand Central Dispatch` instead of `pthreads` , and etc...~~
    - ~~See https://github.com/AntonioMeireles/xhyve/tree/edgy~~
    - Separated [xhyve-bindings](https://github.com/zchee/xhyve-bindings/tree/daemonize)

- [ ] Occasionally fail convert UUID to IP 

- [ ] Support(Ensure) `kill`, `ls`, `restart`, `status`, `stop` command

- [ ] Crash on boot because of `prltoolsd`
    - Crash it's not an empty disk.img?
    - See https://github.com/ailispaw/boot2docker-xhyve/pull/16

- [ ] Cleanup code and more performance
