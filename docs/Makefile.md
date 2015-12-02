# Makefile for Go

# Global go command environment variables
```
ifeq ($(shell [ -d "./Godeps" ] && echo '1' || echo '0'),1)
	GO_CMD := godep go
else
	GO_CMD := go
endif
```
If exist `Godeps` directory, use `godep go` command in all jobs.
```
GO_BUILD=${GO_CMD} build -o ${OUTPUT}
GO_BUILD_RACE=${GO_CMD} build -race -o ${OUTPUT}
GO_TEST=${GO_CMD} test
GO_TEST_VERBOSE=${GO_CMD} test -v
GO_RUN=${GO_CMD} run
GO_INSTALL=${GO_CMD} install -v
GO_CLEAN=${GO_CMD} clean
GO_DEPS=${GO_CMD} get -d -v
GO_DEPS_UPDATE=${GO_CMD} get -d -v -u
GO_VET=${GO_CMD} vet
GO_LINT=golint
```

# Initialized build flags
```
GO_LDFLAGS :=
CGO_ENABLED := 1
```
docker-machine-driver-xhyve use vmnet.framework  
It is binding from C-land to Go
```
CGO_CFLAGS :=
CGO_LDFLAGS :=
CGO_CFLAGS :=
CGO_CPPFLAGS :=
CGO_CXXFLAGS :=
CGO_LDFLAGS :=
GODEBUG :=
```
See https://godoc.org/runtime
```
GOGC :=
```
`GOGC=off go run x.go` or `runtime.SetGCPercent(-1)`  
`-1` for off, `50` for aggressive GC, `100` for default, `200` for lazy GC

# Set debug gcflag, or optimize ldflags
Usage: GDBDEBUG=1 make
```
ifeq ($(GDBDEBUG),1)
	GO_GCFLAGS := -gcflags "-N -l"
```
Disable function inlining and variable registerization. For lldb, gdb, dlv and the involved debugger tools  
See also Dave cheney's blog post: http://goo.gl/6QCJMj  
And, My cgo blog post: http://libraryofalexandria.io/cgo/

`-gcflags '-N'`

- Will be disable the optimisation pass in the compiler  

`-gcflags '-l'`

- Will be disable inlining (but still retain other compiler optimisations) This is very useful if you are investigating small methods, but can’t find them in `objdump`

```
else
	GO_LDFLAGS := $(GO_LDFLAGS) -w -s
```
Turn of DWARF debugging information and strip the binary otherwise  
It will reduce the as much as possible size of the binary  
See also Russ Cox's answered in StackOverflow: http://goo.gl/vOaigc

`-ldflags '-w'`: Turns off DWARF debugging infomation

- Will not be able to use lldb, gdb, objdump or related to debugger tools

`-ldflags '-s'`: Turns off generation of the Go symbol table

- Will not be able to use `go tool nm` to list symbols in the binary
- `strip -s` is like passing '-s' flag to -ldflags, but it doesn't strip quite as much
```
endif
```

## Set static build option
Usage: STATIC=1 make
```
ifeq ($(STATIC),1)
	GO_LDFLAGS := $(GO_LDFLAGS) -extldflags -static
endif
```

## Parse git current branch commit-hash
```
GO_LDFLAGS := ${GO_LDFLAGS} -X `go list ./version`.GitCommit=`git rev-parse --short HEAD 2>/dev/null`
```


# Environment variables

```
export GOARCH=amd64
export GOOS=darwin
```
Hypervisor.framework also vmnet.framework need OS X 10.10 (Yosemite).
See also: https://developer.apple.com/library/mac/releasenotes/MacOSX/WhatsNewInOSX/Articles/MacOSX10_10.html

## Support go1.5 vendoring (let us avoid messing with GOPATH or using godep)
```
export GO15VENDOREXPERIMENT=1
```

## Whether the linker should use external linking mode
```
export GO_EXTLINK_ENABLED=
```
when using -linkmode=auto with code that uses cgo.  
Set to 0 to disable external linking mode, 1 to enable it.


# Package side settings

## Build package infomation
```
GITHUB_USER := zchee
TOP_PACKAGE_DIR := github.com/${GITHUB_USER}
PACKAGE := `basename $(PWD)`
OUTPUT := bin/docker-machine-driver-xhyve
# Parse "func main()" only '.go' file on current dir
# FIXME: Not support main.go
MAIN_FILE := `grep "func main\(\)" *.go -l`
```

## Issue of no include header file in /usr/local/include
See https://github.com/zchee/docker-machine-driver-xhyve/issues/4
```
CGO_CFLAGS=${CGO_CFLAGS} -I/usr/local/include
CGO_LDFLAGS=${CGO_LDFLAGS} -L/usr/local/lib
```

## Include driver debug makefile if `$MACHINE_DEBUG_DRIVER=1`
```
ifeq ($(MACHINE_DEBUG_DRIVER),1)
	include mk/driver.mk
endif
```

# Colorable output
```
CRESET := \x1b[0m
CBLACK := \x1b[30;01m
CRED := \x1b[31;01m
CGREEN := \x1b[32;01m
CYELLOW := \x1b[33;01m
CBLUE := \x1b[34;01m
CMAGENTA := \x1b[35;01m
CCYAN := \x1b[36;01m
CWHITE := \x1b[37;01m
```


# Build jobs settings
```
default: build

makefile-debug:
	@echo ${GO_CMD}

clean:
	@${RM} ./bin/docker-machine-driver-xhyve
	@${RM} ${GOPATH}/bin/docker-machine-driver-xhyve

bin/docker-machine-driver-xhyve:
	sudo chown root:wheel ./bin/docker-machine-driver-xhyve
	sudo chmod u+s ./bin/docker-machine-driver-xhyve

build:
	@echo "${CBLUE}==>${CRESET} Build ${CGREEN}${PACKAGE}${CRESET} ..."
	@echo "${CBLACK} ${GO_BUILD} -ldflags ${GO_LDFLAGS} ${GO_GCFLAGS} ${TOP_PACKAGE_DIR}/${PACKAGE}/bin ${CRESET}"; \
	${GO_BUILD} -ldflags "${GO_LDFLAGS}" ${GO_GCFLAGS} ${TOP_PACKAGE_DIR}/${PACKAGE}/bin || exit 1
	@echo "${CBLUE}==>${CRESET} Change ${CGREEN}${PACKAGE}${CRESET} binary owner and group to root:wheel${CRESET}"; \
	sudo chown root:wheel ${OUTPUT} && sudo chmod u+s ${OUTPUT}


install: bin/docker-machine-driver-xhyve
	sudo cp -p ./bin/docker-machine-driver-xhyve ${GOPATH}/bin/

dep-save:
	godep save $(shell go list ./... | grep -v vendor/)

dep-restore:
	godep restore -v

.PHONY: clean
```
