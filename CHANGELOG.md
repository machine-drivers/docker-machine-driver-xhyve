# [v0.2.0](https://github.com/zchee/docker-machine-driver-xhyve/tree/v0.2.0) (2016-01-09)
[Full Changelog](https://github.com/zchee/docker-machine-driver-xhyve/compare/v0.1.0...v0.2.0)

## Features

### Re-embeded xhyve
Embeded xhyve `C` source again. Not need original xhyve binary also [hooklift/xhyve](https://github.com/zchee/xhyve-bindings).
Without all dependency. It works in `docker-machine-driver-xhyve` one binary.

Use Go bindings for xhyve [hooklift/xhyve](https://github.com/hooklift/xhyve).
And xhyve upsteam repository is [xhyve-xyz/xhyve](https://github.com/xhyve-xyz/xhyve) for now.
Thanks [@johanneswuerbach](https://github.com/johanneswuerbach) and [xhyve-xyz/xhyve](https://github.com/xhyve-xyz/xhyve), [hooklift/xhyve](https://github.com/hooklift/xhyve) developers.

### Sparse (dynamic allocates) volume supoprt
Use OS X sparsebundle to store VM data.
This speeds up volume generation and allocates disk space only on demand.

Thanks [@johanneswuerbach](https://github.com/johanneswuerbach)

### Use a syscall signal to the control of VM
xhyve-xyz was unofficial support `.pid` file.
Now possible to get the state and stop by sending a signal to the process.
It is better than to send `exit 0` use `ssh`.

### Support more than 3GB memory size
[mist64/xhyve](https://github.com/mist64/xhyve) official supported more than 3GB memory size.
See https://github.com/mist64/xhyve/commit/793d17ccffa9a1f74f6f1a4997e73cb2e1496296 .


### Misc
- Rename `--xhyve-memory` flag to `--xhyve-memory-size`. Thanks [@jgeiger](https://github.com/jgeiger)
- Generate `UUID` to use cgo `<uuid.h>` instead of `uuidgen` binary.
- NFS shared folder the more safely use the [johanneswuerbach/nfsexports](https://github.com/johanneswuerbach/nfsexports) package. Thanks [@johanneswuerbach](https://github.com/johanneswuerbach)


## Auto generated Change Log

**Fixed bugs:**

- Can not open boot2docker.iso use hdiutil [\#43](https://github.com/zchee/docker-machine-driver-xhyve/issues/43)

**Merged pull requests:**

- Re-Bump version to 0.2.0 [\#52](https://github.com/zchee/docker-machine-driver-xhyve/pull/52) ([zchee](https://github.com/zchee))
- Fix binary owner to travis:staff and Add multiple Go version release [\#51](https://github.com/zchee/docker-machine-driver-xhyve/pull/51) ([zchee](https://github.com/zchee))
- Back to version 0.1.0 [\#50](https://github.com/zchee/docker-machine-driver-xhyve/pull/50) ([zchee](https://github.com/zchee))
- Add skip\_cleanup [\#48](https://github.com/zchee/docker-machine-driver-xhyve/pull/48) ([zchee](https://github.com/zchee))
- Bump version to 0.2.0 [\#47](https://github.com/zchee/docker-machine-driver-xhyve/pull/47) ([zchee](https://github.com/zchee))
- Update vedor for hooklift/xhyve and change upstream to xhyve-xyz/xhyve [\#46](https://github.com/zchee/docker-machine-driver-xhyve/pull/46) ([zchee](https://github.com/zchee))
- Sparse volume [\#44](https://github.com/zchee/docker-machine-driver-xhyve/pull/44) ([johanneswuerbach](https://github.com/johanneswuerbach))
- Fix Stop\(\) and Kill\(\) to syscall \(SIGTERM|SIGKILL\) [\#42](https://github.com/zchee/docker-machine-driver-xhyve/pull/42) ([zchee](https://github.com/zchee))
- Fix GetState to syscall.Signal use pid [\#41](https://github.com/zchee/docker-machine-driver-xhyve/pull/41) ([zchee](https://github.com/zchee))
- Patched Add -F \<pidfile\> flag to write manage a pidfile [\#39](https://github.com/zchee/docker-machine-driver-xhyve/pull/39) ([zchee](https://github.com/zchee))
- Use xhyve to get the mac address [\#38](https://github.com/zchee/docker-machine-driver-xhyve/pull/38) ([johanneswuerbach](https://github.com/johanneswuerbach))
- Actually support CPU -1 for all available CPUs [\#36](https://github.com/zchee/docker-machine-driver-xhyve/pull/36) ([johanneswuerbach](https://github.com/johanneswuerbach))
- Auto-create missing exports file, line break fixes [\#35](https://github.com/zchee/docker-machine-driver-xhyve/pull/35) ([johanneswuerbach](https://github.com/johanneswuerbach))
- Fix travis-ci settings [\#34](https://github.com/zchee/docker-machine-driver-xhyve/pull/34) ([zchee](https://github.com/zchee))
- Enhanced NFS shares [\#33](https://github.com/zchee/docker-machine-driver-xhyve/pull/33) ([johanneswuerbach](https://github.com/johanneswuerbach))
- Fixes xhyve stdout and stderr printing [\#32](https://github.com/zchee/docker-machine-driver-xhyve/pull/32) ([johanneswuerbach](https://github.com/johanneswuerbach))
- Improved boot2docker version detection [\#31](https://github.com/zchee/docker-machine-driver-xhyve/pull/31) ([johanneswuerbach](https://github.com/johanneswuerbach))
- Fixed xhyve-memory-size [\#30](https://github.com/zchee/docker-machine-driver-xhyve/pull/30) ([johanneswuerbach](https://github.com/johanneswuerbach))
- Use travis with OS X 10.11.1 [\#28](https://github.com/zchee/docker-machine-driver-xhyve/pull/28) ([johanneswuerbach](https://github.com/johanneswuerbach))
- Embeds xhyve [\#27](https://github.com/zchee/docker-machine-driver-xhyve/pull/27) ([johanneswuerbach](https://github.com/johanneswuerbach))
- Use --xhyve-memory-size to match other docker-machine drivers [\#26](https://github.com/zchee/docker-machine-driver-xhyve/pull/26) ([jgeiger](https://github.com/jgeiger))
- Update README for memory command [\#25](https://github.com/zchee/docker-machine-driver-xhyve/pull/25) ([jgeiger](https://github.com/jgeiger))
- Bump v0.1.0 [\#23](https://github.com/zchee/docker-machine-driver-xhyve/pull/23) ([zchee](https://github.com/zchee))


# [v0.1.0](https://github.com/zchee/docker-machine-driver-xhyve/tree/v0.1.0) (2015-12-02)
[Full Changelog](https://github.com/zchee/docker-machine-driver-xhyve/compare/v0.0.1...v0.1.0)

## Auto generated Change Log

**Closed issues:**

- Shared folders [\#20](https://github.com/zchee/docker-machine-driver-xhyve/issues/20)
- unrecognized import path "libguestfs.org/guestfs" [\#17](https://github.com/zchee/docker-machine-driver-xhyve/issues/17)
- Occasionally, hangs at get ip [\#15](https://github.com/zchee/docker-machine-driver-xhyve/issues/15)
- godep: command not found [\#14](https://github.com/zchee/docker-machine-driver-xhyve/issues/14)
- Command not found: VBoxManage [\#13](https://github.com/zchee/docker-machine-driver-xhyve/issues/13)
- docker-machine rm hangs [\#8](https://github.com/zchee/docker-machine-driver-xhyve/issues/8)
- No vendored dependencies [\#7](https://github.com/zchee/docker-machine-driver-xhyve/issues/7)
- 'guestfs.h' file not found [\#4](https://github.com/zchee/docker-machine-driver-xhyve/issues/4)
- unrecognized import path "libguestfs.org/guestfs" [\#2](https://github.com/zchee/docker-machine-driver-xhyve/issues/2)

**Merged pull requests:**

- Experimental: Auto-create NFS share [\#22](https://github.com/zchee/docker-machine-driver-xhyve/pull/22) ([johanneswuerbach](https://github.com/johanneswuerbach))
- Fix makefile env [\#21](https://github.com/zchee/docker-machine-driver-xhyve/pull/21) ([zchee](https://github.com/zchee))
- Generate disk image using hdiutil [\#19](https://github.com/zchee/docker-machine-driver-xhyve/pull/19) ([johanneswuerbach](https://github.com/johanneswuerbach))
- Trimming only "0" of the ten's digit [\#16](https://github.com/zchee/docker-machine-driver-xhyve/pull/16) ([zchee](https://github.com/zchee))
- Remove hardcode boot2docker version, will parse [\#12](https://github.com/zchee/docker-machine-driver-xhyve/pull/12) ([zchee](https://github.com/zchee))
- Implement all of the docker-machine commands [\#11](https://github.com/zchee/docker-machine-driver-xhyve/pull/11) ([zchee](https://github.com/zchee))
- Fix GetIP and GetState to Correctly get vm status [\#10](https://github.com/zchee/docker-machine-driver-xhyve/pull/10) ([zchee](https://github.com/zchee))
- Update vendor/github.com/docker/machine subtree [\#5](https://github.com/zchee/docker-machine-driver-xhyve/pull/5) ([pierrezurek](https://github.com/pierrezurek))


# [v0.0.1](https://github.com/zchee/docker-machine-driver-xhyve/tree/v0.0.1) (2015-09-27)

## Auto generated Change Log

**Merged pull requests:**

- Behavior boot xhyve - carry to my repo from nathanleclaire/docker-machine-xhyve [\#1](https://github.com/zchee/docker-machine-driver-xhyve/pull/1) ([zchee](https://github.com/zchee))
