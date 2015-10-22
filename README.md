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
Or, See experimental embedded xhyve branch [embed-xhyve](https://github.com/zchee/docker-machine-xhyve/tree/embed-xhyve)

```bash
$ go get -d github.com/zchee/xhyve-bindings
$ cd $GOPATH/src/github.com/zchee/xhyve-bindings
$ git checkout daemonize
$ make
$ make install
```

### libguestfs
http://libguestfs.org/

Create `ext.4` filesystem disk image using libguestfs golang binding.

```bash
$ brew tap zchee/libguestfs
$ brew install libguestfs --with-go --devel --env=std
```
Current status only support golang binding.

Also, downloading `supermin appliance` kernel files.  
Warning! Kernel file size over 4GB!


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

- [x] ~~Occasionally fail convert UUID to IP~~
    - Fixed [1960629b3c8683aec193631a0e9573c5143832ab](https://github.com/zchee/docker-machine-xhyve/commit/1960629b3c8683aec193631a0e9573c5143832ab)

- [ ] Support(Ensure) `kill`, `ls`, `restart`, `status`, `stop` command

- [ ] Crash on boot because of `prltoolsd`
    - Crash it's not an empty disk.img?
    - See https://github.com/ailispaw/boot2docker-xhyve/pull/16

- [ ] Cleanup code and more performance
