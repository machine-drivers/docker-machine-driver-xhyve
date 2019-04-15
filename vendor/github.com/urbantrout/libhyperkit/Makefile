# -----------------------------------------------------------------------------
# includeing xhyve original config.mk
-include xhyve.mk


# -----------------------------------------------------------------------------
# ocaml-qcow bindings

HAVE_OCAML_QCOW := $(shell if ocamlfind query qcow uri >/dev/null 2>/dev/null ; then echo YES ; else echo NO; fi)

ifeq ($(HAVE_OCAML_QCOW),YES)
CGO_CFLAGS += -DHAVE_OCAML=1 -DHAVE_OCAML_QCOW=1 -DHAVE_OCAML=1

OCAML_WHERE := $(shell ocamlc -where)
OCAML_LDLIBS := -L$(OCAML_WHERE) \
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

build: CGO_CFLAGS += -I$(OCAML_WHERE)
build: CGO_LDFLAGS += $(OCAML_LDLIBS)
build: GO_BUILD_TAGS += qcow2
build: generate
endif


# -----------------------------------------------------------------------------
# make rules

build:
	CGO_CFLAGS="$(CGO_CFLAGS)" CGO_LDFLAGS="$(CGO_LDFLAGS)" go build -v -x -tags=$(GO_BUILD_TAGS) .

mirage_block_ocaml.o:
	go generate -v -x -tags=$(GO_BUILD_TAGS)

generate: mirage_block_ocaml.o


vendor-hyperkit:
	-git clone https://github.com/docker/hyperkit.git hyperkit

patch-generate: patch-apply
	-cd hyperkit; git diff > ../xhyve.patch

patch-apply: vendor-hyperkit
	cd hyperkit; \
		for p in $(shell find patch \( -name \*.patch \)); do \
			patch -Nl -p1 -F4 < ../$$p; \
		done

sync: clean vendor-hyperkit apply-patch
	find . \( -name \*.orig -o -name \*.rej \) -delete
	for file in $(SRC); do \
		cp -f $$file $$(basename $$file) ; \
	done
	cp -r hyperkit/src/include include
	cp hyperkit/README.md README.hyperkit.md
	cp hyperkit/README.xhyve.md .


clean:
	${RM} *.a *.o *.cmi *.cmx
	${RM} -r hyperkit

.PHONY: build generate vendor-hyperkit patch-generate patch-apply sync clean
