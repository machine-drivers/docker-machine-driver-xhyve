# Change Log

## [v0.3.0](https://github.com/zchee/docker-machine-driver-xhyve/tree/v0.3.0) (2016-11-06)
[Full Changelog](https://github.com/zchee/docker-machine-driver-xhyve/compare/v0.2.3...v0.3.0)

## Features

- Suuport rkt container engine for minikube. See [kubernetes/minikube#using-rkt-container-engine](https://github.com/kubernetes/minikube/#using-rkt-container-engine)
- Change [hooklift/xhyve](https://github.com/hooklift/xhyve) to [zchee/lihyperkit](https://github.com/zchee/lihyperkit)
- Support QCow2 disk image using [docker/hyperkit](https://github.com/docker/hyperkit)'s [ocaml-qcow](https://github.com/mirage/ocaml-qcow) bindings
- Supprot configurable kernel and initrd filename

## Auto generated Change Log

**Closed issues:**

- Make /boot configurable to allow different live iso to work with driver  [\#146](https://github.com/zchee/docker-machine-driver-xhyve/issues/146)

**Merged pull requests:**

- Bump version to 0.3.0 [\#149](https://github.com/zchee/docker-machine-driver-xhyve/pull/149) ([zchee](https://github.com/zchee))
- Fix \#146 Make /boot configurable to allow different live iso to work with driver [\#147](https://github.com/zchee/docker-machine-driver-xhyve/pull/147) ([praveenkumar](https://github.com/praveenkumar))
- ci/travis: add setup opam and libev for ocaml-qcow [\#145](https://github.com/zchee/docker-machine-driver-xhyve/pull/145) ([zchee](https://github.com/zchee))
- Auto-detect whether the kernel is bzImage or vmlinuz. [\#144](https://github.com/zchee/docker-machine-driver-xhyve/pull/144) ([dlorenc](https://github.com/dlorenc))
- Defer the unmount. [\#143](https://github.com/zchee/docker-machine-driver-xhyve/pull/143) ([dlorenc](https://github.com/dlorenc))
- Add some debug logging during mac address retrieval. [\#142](https://github.com/zchee/docker-machine-driver-xhyve/pull/142) ([dlorenc](https://github.com/dlorenc))
- Support the rkt/minikube-iso [\#140](https://github.com/zchee/docker-machine-driver-xhyve/pull/140) ([zchee](https://github.com/zchee))
- ci/travis: upgdate go version to 1.7.1 and remove 1.6.3 & fix install [\#139](https://github.com/zchee/docker-machine-driver-xhyve/pull/139) ([zchee](https://github.com/zchee))
- xhyve/qcow2: support qcow2 disk image format [\#138](https://github.com/zchee/docker-machine-driver-xhyve/pull/138) ([zchee](https://github.com/zchee))
- Support using the configured UUID instead of always generating a UUID… [\#133](https://github.com/zchee/docker-machine-driver-xhyve/pull/133) ([chirino](https://github.com/chirino))
- CI: upgrade TravisCI go version to 1.6.3 and 1.7rc5 [\#131](https://github.com/zchee/docker-machine-driver-xhyve/pull/131) ([zchee](https://github.com/zchee))

## [v0.2.3](https://github.com/zchee/docker-machine-driver-xhyve/tree/v0.2.3) (2016-08-05)
[Full Changelog](https://github.com/zchee/docker-machine-driver-xhyve/compare/v0.2.2...v0.2.3)

**Closed issues:**

- Please declare a license [\#108](https://github.com/zchee/docker-machine-driver-xhyve/issues/108)
- nfs share permissions problem [\#99](https://github.com/zchee/docker-machine-driver-xhyve/issues/99)
- Build issue with docker v1.10 [\#92](https://github.com/zchee/docker-machine-driver-xhyve/issues/92)
- Could not create vmnet interface, permission denied or no entitlement? [\#85](https://github.com/zchee/docker-machine-driver-xhyve/issues/85)
- B2D iso and directory created by root owner [\#82](https://github.com/zchee/docker-machine-driver-xhyve/issues/82)

**Merged pull requests:**

- Bump version to 0.2.3 [\#130](https://github.com/zchee/docker-machine-driver-xhyve/pull/130) ([zchee](https://github.com/zchee))
- Changed doc to reflect that virtio-9p support is included in the newer releases of boot2docker [\#129](https://github.com/zchee/docker-machine-driver-xhyve/pull/129) ([r2d4](https://github.com/r2d4))
- Add error checking to CopyIsoToMachineDir. [\#127](https://github.com/zchee/docker-machine-driver-xhyve/pull/127) ([dlorenc](https://github.com/dlorenc))
- Fix lib9p build [\#106](https://github.com/zchee/docker-machine-driver-xhyve/pull/106) ([zchee](https://github.com/zchee))
- Set the NFS export to map all the users on the remote to the current user [\#105](https://github.com/zchee/docker-machine-driver-xhyve/pull/105) ([jamesRaybould](https://github.com/jamesRaybould))
- fix check for binary root ownership [\#104](https://github.com/zchee/docker-machine-driver-xhyve/pull/104) ([codekitchen](https://github.com/codekitchen))
- Support virtio-9p [\#95](https://github.com/zchee/docker-machine-driver-xhyve/pull/95) ([zchee](https://github.com/zchee))
- Update travis to xcode 7.3 beta [\#89](https://github.com/zchee/docker-machine-driver-xhyve/pull/89) ([zchee](https://github.com/zchee))
- Update vendor for docker/machine [\#88](https://github.com/zchee/docker-machine-driver-xhyve/pull/88) ([zchee](https://github.com/zchee))
- Self implement b2dutils.CopyIsoToMachineDir [\#87](https://github.com/zchee/docker-machine-driver-xhyve/pull/87) ([zchee](https://github.com/zchee))

# [v0.2.2](https://github.com/zchee/docker-machine-driver-xhyve/tree/v0.2.2) (2016-01-22)
[Full Changelog](https://github.com/zchee/docker-machine-driver-xhyve/compare/v0.2.1...v0.2.2)

## Features

### Check own binary owner and execute user
Added PreCommandCheck method for required of docker-machine-driver-xhyve.

### Misc
- Fixed GetUrl does not handle GetState

## Auto generated Change Log

**Fixed bugs:**

- Check owner before every command [\#77](https://github.com/zchee/docker-machine-driver-xhyve/issues/77)
- Does not handle sudo before docker-machine commands [\#76](https://github.com/zchee/docker-machine-driver-xhyve/issues/76)
- GetUrl\(\) does not handle GetState\(\) [\#68](https://github.com/zchee/docker-machine-driver-xhyve/issues/68)

**Closed issues:**

- Using docker-machine-driver-xhyve without sudo privledges [\#74](https://github.com/zchee/docker-machine-driver-xhyve/issues/74)
- Show a nice warning when running without root [\#64](https://github.com/zchee/docker-machine-driver-xhyve/issues/64)
- about IP address [\#62](https://github.com/zchee/docker-machine-driver-xhyve/issues/62)
- ip not found in dhcp leases [\#60](https://github.com/zchee/docker-machine-driver-xhyve/issues/60)
- invalid flag argument in current master [\#59](https://github.com/zchee/docker-machine-driver-xhyve/issues/59)
- Any plan about HomeBrew formulae support? [\#24](https://github.com/zchee/docker-machine-driver-xhyve/issues/24)

**Merged pull requests:**

- Bump version to 0.2.2 [\#81](https://github.com/zchee/docker-machine-driver-xhyve/pull/81) ([zchee](https://github.com/zchee))
- Update and cleanup vendor [\#80](https://github.com/zchee/docker-machine-driver-xhyve/pull/80) ([zchee](https://github.com/zchee))
- Add PreCommandCheck [\#78](https://github.com/zchee/docker-machine-driver-xhyve/pull/78) ([zchee](https://github.com/zchee))
- Add check binary owner on PreCreateCheck\(\) [\#75](https://github.com/zchee/docker-machine-driver-xhyve/pull/75) ([zchee](https://github.com/zchee))
- Rename test vm name, Remove kill job, Change brackets [\#73](https://github.com/zchee/docker-machine-driver-xhyve/pull/73) ([zchee](https://github.com/zchee))
- Fix handling vm status before GetURL\(\) [\#72](https://github.com/zchee/docker-machine-driver-xhyve/pull/72) ([zchee](https://github.com/zchee))
- Change go version to 1.5.3 and 1.6beta2 [\#66](https://github.com/zchee/docker-machine-driver-xhyve/pull/66) ([zchee](https://github.com/zchee))
- Add circleci badge and markdown table [\#65](https://github.com/zchee/docker-machine-driver-xhyve/pull/65) ([zchee](https://github.com/zchee))
- \[WIP\] Add CircleCI use osxcc [\#63](https://github.com/zchee/docker-machine-driver-xhyve/pull/63) ([zchee](https://github.com/zchee))
- ci: Fix multiple Go version release [\#61](https://github.com/zchee/docker-machine-driver-xhyve/pull/61) ([zchee](https://github.com/zchee))

# [v0.2.1](https://github.com/zchee/docker-machine-driver-xhyve/tree/v0.2.1) (2016-01-12)
[Full Changelog](https://github.com/zchee/docker-machine-driver-xhyve/compare/v0.2.0...v0.2.1)

## Features

### Add wait for SSH
Add wait for available SSH login when start and restart commands.

### Misc
- Fix folder structure. Thanks [@saljam](https://github.com/saljam).

## Auto generated Change Log

**Closed issues:**

- make docker-machine-driver-xhyve go gettable [\#49](https://github.com/zchee/docker-machine-driver-xhyve/issues/49)
- Might want to re-release a binary [\#45](https://github.com/zchee/docker-machine-driver-xhyve/issues/45)
- Build isn't including the updated --memory-size flags? [\#40](https://github.com/zchee/docker-machine-driver-xhyve/issues/40)

**Merged pull requests:**

- Bump version to 0.2.1 [\#58](https://github.com/zchee/docker-machine-driver-xhyve/pull/58) ([zchee](https://github.com/zchee))
- Add wait for SSH when \(Re\)Start [\#57](https://github.com/zchee/docker-machine-driver-xhyve/pull/57) ([zchee](https://github.com/zchee))
- Update vendor for docker/docker and docker/machine [\#56](https://github.com/zchee/docker-machine-driver-xhyve/pull/56) ([zchee](https://github.com/zchee))
- Fix Makefile for develop [\#55](https://github.com/zchee/docker-machine-driver-xhyve/pull/55) ([zchee](https://github.com/zchee))
- make main binary go installable [\#53](https://github.com/zchee/docker-machine-driver-xhyve/pull/53) ([saljam](https://github.com/saljam))

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
