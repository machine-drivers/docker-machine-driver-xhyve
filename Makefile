#
#  Makefile for Go
#


# ----------------------------------------------------------------------------
# debug flag
ifeq ($V, 1)
	VERBOSE =
	GO_VERBOSE = -v -x
else
	VERBOSE = @
	GO_VERBOSE =
endif


# ----------------------------------------------------------------------------
# Include makefiles

-include mk/color.mk
-include mk/lib9p.mk
# Include driver debug makefile if $MACHINE_DRIVER_DEBUG=1
ifeq ($(MACHINE_DRIVER_DEBUG),1)
-include mk/driver.mk
endif


# ----------------------------------------------------------------------------
# Package settings

# Build package infomation
GITHUB_USER := zchee
TOP_PACKAGE_DIR := github.com/${GITHUB_USER}
PACKAGE := $(shell basename $(PWD))
OUTPUT := bin/docker-machine-driver-xhyve
# Parse "func main()" only '.go' file on current dir
# FIXME: Not support main.go
MAIN_FILE := $(shell grep "func main\(\)" *.go -l)


# ----------------------------------------------------------------------------
# Define main commands

CC := $(shell xcrun -f clang)
LIBTOOL := $(shell xcrun -f libtool)
GO_CMD := $(shell which go)
GIT_CMD := $(shell which git)
DOCKER_CMD := $(shell which docker)


# ----------------------------------------------------------------------------
# Define sub commmands

GO_BUILD = ${GO_CMD} build $(GO_VERBOSE) $(GO_BUILD_FLAG) -o ${OUTPUT}

GO_INSTALL = ${GO_CMD} install $(GO_BUILD_TAGS)

GO_RUN = ${GO_CMD} run

GO_TEST = ${GO_CMD} test ${GO_VERBOSE}
GO_TEST_RUN = ${GO_TEST} -run ${RUN}
GO_TEST_ALL = test -race -cover -bench=.
GO_VET = ${GO_CMD} vet
GO_LINT = golint

GO_DEPS=${GO_CMD} get -d
GO_DEPS_UPDATE=${GO_CMD} get -d -u
# Check godep binary
GODEP = ${GOPATH}/bin/godep
GODEP_CMD = $(if ${GODEP}, , $(error Please install godep: go get github.com/tools/godep)) ${GODEP}


# ----------------------------------------------------------------------------
# Define build flags

# Initialize
CGO_CFLAGS :=
CGO_LDFLAGS :=

# Parse git current branch commit-hash
GO_LDFLAGS += -X `go list ./xhyve`.GitCommit=`git rev-parse --short HEAD 2>/dev/null`


# Set debug gcflag, or optimize ldflags
#  Usage: DEBUG=true make
ifeq ($(DEBUG),true)
GO_GCFLAGS += -N -l
# Disable function inlining and variable registerization. For lldb, gdb, dlv and the involved debugger tools
# See also Dave cheney's blog post: http://goo.gl/6QCJMj
# And, My cgo blog post: http://libraryofalexandria.io/cgo/
#
# -gcflags '-N': Will be disable the optimisation pass in the compiler
# -gcflags '-l': Will be disable inlining (but still retain other compiler optimisations)
#                This is very useful if you are investigating small methods, but can’t find them in `objdump`
else
GO_LDFLAGS += -w -s
# Turn of DWARF debugging information and strip the binary otherwise
# It will reduce the as much as possible size of the binary
# See also Russ Cox's answered in StackOverflow: http://goo.gl/vOaigc
#
# -ldflags '-w': Turns off DWARF debugging infomation
# 	- Will not be able to use lldb, gdb, objdump or related to debugger tools
# -ldflags '-s': Turns off generation of the Go symbol table
# 	- Will not be able to use `go tool nm` to list symbols in the binary
# 	- `strip -s` is like passing '-s' flag to -ldflags, but it doesn't strip quite as much
endif


# build with race
ifeq ($(RACE),true)
GO_BUILD_FLAG += -race
endif


# Use virtio-9p shared folder with static build lib9p library
# Default is enable
GO_BUILD_TAGS ?= lib9p
# included 'lib9p' in the $GO_BUILD_TAGS, and exists 'lib9p.a' file
ifneq (,$(findstring lib9p,$(GO_BUILD_TAGS)))
CGO_CFLAGS += -I${PWD}/vendor/lib9p
CGO_LDFLAGS += ${PWD}/vendor/build/lib9p/lib9p.a
bin/docker-machine-driver-xhyve: lib9p
endif


# Use mirage-block for pwritev|preadv
HAVE_OCAML_QCOW := $(shell if ocamlfind query qcow uri >/dev/null 2>/dev/null ; then echo YES ; else echo NO; fi)
# included 'qcow2' in the $GO_BUILD_TAGS, and $HAVE_OCAML_QCOW is YES
ifneq (,$(findstring qcow2,$(GO_BUILD_TAGS)))
ifeq ($(HAVE_OCAML_QCOW),YES)
OCAML_WHERE := $(shell ocamlc -where)
OCAML_LDLIBS := -L $(OCAML_WHERE) \
	$(shell ocamlfind query cstruct)/cstruct.a \
	$(shell ocamlfind query cstruct)/libcstruct_stubs.a \
	$(shell ocamlfind query io-page)/io_page.a \
	$(shell ocamlfind query io-page)/io_page_unix.a \
	$(shell ocamlfind query io-page)/libio_page_unix_stubs.a \
	$(shell ocamlfind query lwt.unix)/liblwt-unix_stubs.a \
	$(shell ocamlfind query lwt.unix)/lwt-unix.a \
	$(shell ocamlfind query lwt.unix)/lwt.a \
	$(shell ocamlfind query threads)/libthreadsnat.a \
	$(shell ocamlfind query mirage-block-unix)/libmirage_block_unix_stubs.a \
	-lasmrun -lbigarray -lunix
CGO_CFLAGS += -DHAVE_OCAML=1 -DHAVE_OCAML_QCOW=1 -DHAVE_OCAML=1 -I$(OCAML_WHERE)
CGO_LDFLAGS += $(OCAML_LDLIBS)
bin/docker-machine-driver-xhyve: vendor/github.com/zchee/libhyperkit/mirage_block_ocaml.syso
endif
endif


GO_BUILD_FLAG += -tags='$(GO_BUILD_TAGS)'


# ----------------------------------------------------------------------------
# Go environment variables

# Hypervisor.framework also vmnet.framework need OS X 10.10 (Yosemite).
# See also:
#   https://developer.apple.com/library/mac/releasenotes/MacOSX/WhatsNewInOSX/Articles/MacOSX10_10.html
export GOARCH=amd64
export GOOS=darwin

# Support go1.5 vendoring (let us avoid messing with GOPATH or using godep)
export GO15VENDOREXPERIMENT=1

# TODO: uuid.go need cgo
export CGO_ENABLED=1


# ----------------------------------------------------------------------------
# Build jobs settings

default: build

build: bin/docker-machine-driver-xhyve

vendor/github.com/zchee/libhyperkit/mirage_block_ocaml.syso:
	$(VERBOSE) $(GO_CMD) generate $(GO_BUILD_FLAG) $(GO_VERBOSE) ./vendor/github.com/zchee/libhyperkit

bin/docker-machine-driver-xhyve:
	$(VERBOSE) test -d bin || mkdir -p bin;
	@echo "${CBLUE}==>${CRESET} Build ${CGREEN}${PACKAGE}${CRESET}..."
	$(VERBOSE) $(ENV) CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" $(GO_BUILD) -gcflags "$(GO_GCFLAGS)" -ldflags "$(GO_LDFLAGS)" ${TOP_PACKAGE_DIR}/${PACKAGE}

build-privilege: bin/docker-machine-driver-xhyve
	@echo "${CBLUE}==>${CRESET} Change ${CGREEN}${PACKAGE}${CRESET} binary owner and group to root:wheel. Please root password${CRESET}"
	$(VERBOSE) $(ENV) sudo chown root:wheel ${OUTPUT} && sudo chmod u+s ${OUTPUT}

install: build-privilege
	@echo "${CBLUE}==>${CRESET} Install ${CGREEN}${PACKAGE}${CRESET}..."
	$(VERBOSE) test -d /usr/local/bin || mkdir -p /usr/local/bin
	sudo cp -p ./bin/docker-machine-driver-xhyve /usr/local/bin/


test:
	@echo "${CBLUE}==>${CRESET} Test ${CGREEN}${PACKAGE}${CRESET}..."
	@echo "${CBLACK} ${GO_TEST} ${TOP_PACKAGE_DIR}/${PACKAGE}/xhyve ${CRESET}"; \
	${GO_TEST} ${TOP_PACKAGE_DIR}/${PACKAGE}/xhyve || exit 1

test-bindings:
	$(VERBOSE) if nm bin/docker-machine-driver-xhyve | grep _l9p_server_init >/dev/null 2>&1; then echo 'lib9p'; fi
	$(VERBOSE) if nm bin/docker-machine-driver-xhyve | grep _camlMirage_block__code_begin >/dev/null 2>&1; then echo 'qcow2'; fi


vendor-update:
	$(VERBOSE) gvt update -all

vendor-restore:
	$(VERBOSE) gvt restore


docker-build:
	${DOCKER_CMD} build --rm -t ${GITHUB_USER}/${PACKAGE} .

docker-build-nocache:
	${DOCKER_CMD} build --rm --no-cache -t ${GITHUB_USER}/${PACKAGE} .


clean: clean-lib9p
	@${RM} -r ./bin ./vendor/github.com/zchee/libhyperkit/*.cmi ./vendor/github.com/zchee/libhyperkit/*.cmx ./vendor/github.com/zchee/libhyperkit/*.syso


run: driver-run

rm: driver-remove

# TODO: for zsh completion. zsh do not get jobs of includes makefile
test-env:
test-inspect:
test-ip:
test-kill:
test-ls:
test-regenerate-certs:
test-restart:
test-rm:
test-ssh:
test-start:
test-status:
test-stop:
test-upgrade:
test-url:

.PHONY: clean run rm kill
