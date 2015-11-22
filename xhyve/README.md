# Xhyve from Go

The goal of this project is to compile [xhyve](https://github.com/mist64/xhyve) into a Go
package and be able to distribute a Go binary with xhyve embedded.

Currently, the bindings are only for the `main()` entrypoint in `xhyve.c`, allowing
the Go program to pass in any command line arguments to xhyve. This is a stop-gap for now,
and I welcome any effort to make actual Go bindings to the underlying xhyve functions.

# Installation

I only tested on OS X Yosemite, which is the first OS X version to have Hypervisor.framework
which is what xhyve leverages.

```shell
$ go get github.com/tiborvass/xhyve-bindings
$ cd $GOPATH/src/github.com/tiborvass/xhyve-bindings
$ go run main/main.go upstream/test/vmlinuz upstream/test/initrd.gz
```

Once in the VM, type `sudo halt` to quit it.

# Documentation

Follow [Godoc](https://godoc.org/github.com/tiborvass/xhyve-bindings).

# Roadmap

* **fork-exec**: Go program should be able to run multiple xhyve instances. I suggest we use
the [reexec package](https://godoc.org/github.com/docker/docker/pkg/reexec).
* **cross-compiling**: We should be able to compile it from Linux. I was thinking of using
[xgo](https://github.com/karalabe/xgo) but any other solution that works should be good.
* **management**: Start, Stop and Kill fork-exec'd xhyve instance.
* **pty**: Not sure how important this is, but I wanted to have a way to attach and detach from
the TTY.

# Contributing

Just send pull requests, open issues.

The `upstream/` directory is a git subtree. It's the first time I use that functionality
so bear with me. Suggestions welcome.

If you need to update upstream, run `make clean` first, update upstream, and run `make`.
It will apply a small patch `upstream.patch` that's currently needed.

By using cgo, we're limited to requiring all `*.c` files at the root of the repository, hence
all those symlinks created by `generate.sh` when doing `make`. This also means, that no two
C files can have the same name.

# License

The bindings are under MIT License.
For xhyve itself, read [https://github.com/mist64/xhyve](https://github.com/mist64/xhyve).
