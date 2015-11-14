docker-machine-xhyve
===

Docker Machine driver plugin for [xhyve](https://github.com/mist64/xhyve) native OS X Hypervisor

Master branch inherited from [nathanleclaire/docker-machine-xhyve](https://github.com/nathanleclaire/docker-machine-xhyve). Thanks [@nathanleclaire](https://github.com/nathanleclaire) :)  
If you have issues or pull-requests, Desired to be posted to this repository.


## Required

### docker-machine
https://github.com/docker/machine

**!! Please do not post the issue of this repository to the docker/machine !!**  
It will interfere with the development of the docker-machine.  
If you were doubt problem either, please post to this repository. I will judge.

Now, `libmachine-rpc` driver plugin method is merged `docker-machine` master branch.  
https://github.com/docker/machine/commit/8aa1572e0dcd75762a7627e1056ef104317f44b9
Awesome @nathanleclaire :tada:

```bash
go get github.com/nathanleclaire/machine
cd $GOPATH/src/github.com/docker/machine
# Build docker-machine and some docker-machine official(embedded) driver binary
make build
# Install all binary into /usr/local/bin/
make install
```

### xhyve-bindings
https://github.com/zchee/xhyve-bindings

Since it is was hard to `os.exec` itself that embedded `xhyve.Exec`, for the time being, is separated into [xhyve-bindings](https://github.com/zchee/xhyve-bindings).  
Or, See experimental embedded xhyve branch [embed-xhyve](https://github.com/zchee/docker-machine-xhyve/tree/embed-xhyve)

```bash
$ go get -d github.com/zchee/xhyve-bindings
$ cd $GOPATH/src/github.com/zchee/xhyve-bindings
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
$ cd $GOPATH/src/github.com/zchee/docker-machine-xhyve
$ make
$ make install
```


## TODO

- [ ] Shared folder support
  - Use `9p` filesystem also `virtio-9p`? See https://github.com/mist64/xhyve/issues/70#issuecomment-144935541

- [ ] Replace execute binary to Go `syscall`
    - [ ] `hdutil`
    - [x]  `dd`
      - Create ext.4 filesystem disk use `libguestfs`

- [ ] Replace generate uuid, native Go code instead of `uuidgen`

- [ ] Cleanup code and more performance
