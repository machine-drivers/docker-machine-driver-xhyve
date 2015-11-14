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


## Known isuue

### Occasionally, hangs at get ip (MACaddress to IP)
The cause is not still unknown, it is continuing to investigate.

### Not implement shared folder
`docker-machine-xhyve` also `xhyve` does not implement shared folder system.  
Please use the existing protocol for the time being, such as `sshd`.

### Get state use `ssh`, do not know real vm state
`docker-machine-xhyve` checking vm state use send `exit 0` on `ssh`.  
but, that is not real state of vm.  

`xhyve` has the XPC dictionary on backend. (still uncertain)  
XPC Services is part of libSystem, provides a lightweight mechanism for basic interprocess communication integrated with Grand Central Dispatch and launchd.  
In the future, get the state to use it.

### Does not clean up the vmnet when remove a VM
Current state, `docker-machine-xhyve` does not clean up the vmnet configuration.  

```
Running xhyve vm (e.g. IP: 192.168.64.1)
        |
`docker-machine rm`, ACPI signal(poweroff) in over ssh
        |
vm is poweroff, `goxhyve` also killed with to poweroff
        |
Re-create xhyve vm.
```
New vm's ip, It will probably be assigned to 192.168.64.**2**. If create another new vm, assigned to 192.168.64.**3**  
but 192.168.64.**1** are not using anyone.

`vmnet.framework` seems to have decided to IP based on `/var/db/dhcpd_leases` and `/Library/Preferences/SystemConfiguration/com.apple.vmnet.plist`  
So, To remove it manually, or Don’t even bother.  
Maybe `vmnet.framework` shared net address range `192.168.64.1` ~ `192.168.64.255`. You can make 255 vm :stuck_out_tongue_closed_eyes:

I will fix after I understand the `vmnet.framework`

## TODO

- [ ] Shared folder support
  - Use `9p` filesystem also `virtio-9p`? See https://github.com/mist64/xhyve/issues/70#issuecomment-144935541

- [ ] Replace execute binary to Go `syscall`
    - [ ] `hdutil`
    - [x]  `dd`
      - Create ext.4 filesystem disk use `libguestfs`

- [ ] Replace generate uuid, native Go code instead of `uuidgen`

- [ ] Cleanup code and more performance
