XHYVE_VERSION := $(shell cd hyperkit; git describe --abbrev=6 --dirty --always --tags)
GIT_VERSION_SHA1 := $(shell cd hyperkit; git rev-parse HEAD)

VMM_SRC := \
	hyperkit/src/vmm/x86.c \
	hyperkit/src/vmm/vmm.c \
	hyperkit/src/vmm/vmm_host.c \
	hyperkit/src/vmm/vmm_mem.c \
	hyperkit/src/vmm/vmm_lapic.c \
	hyperkit/src/vmm/vmm_instruction_emul.c \
	hyperkit/src/vmm/vmm_ioport.c \
	hyperkit/src/vmm/vmm_callout.c \
	hyperkit/src/vmm/vmm_stat.c \
	hyperkit/src/vmm/vmm_util.c \
	hyperkit/src/vmm/vmm_api.c \
	hyperkit/src/vmm/intel/vmx.c \
	hyperkit/src/vmm/intel/vmx_msr.c \
	hyperkit/src/vmm/intel/vmcs.c \
	hyperkit/src/vmm/io/vatpic.c \
	hyperkit/src/vmm/io/vatpit.c \
	hyperkit/src/vmm/io/vhpet.c \
	hyperkit/src/vmm/io/vioapic.c \
	hyperkit/src/vmm/io/vlapic.c \
	hyperkit/src/vmm/io/vpmtmr.c \
	hyperkit/src/vmm/io/vrtc.c

XHYVE_SRC := \
	hyperkit/src/acpitbl.c \
	hyperkit/src/atkbdc.c \
	hyperkit/src/block_if.c \
	hyperkit/src/consport.c \
	hyperkit/src/dbgport.c \
	hyperkit/src/inout.c \
	hyperkit/src/ioapic.c \
	hyperkit/src/md5c.c \
	hyperkit/src/mem.c \
	hyperkit/src/mevent.c \
	hyperkit/src/mptbl.c \
	hyperkit/src/pci_ahci.c \
	hyperkit/src/pci_emul.c \
	hyperkit/src/pci_hostbridge.c \
	hyperkit/src/pci_irq.c \
	hyperkit/src/pci_lpc.c \
	hyperkit/src/pci_uart.c \
	hyperkit/src/pci_virtio_9p.c \
	hyperkit/src/pci_virtio_block.c \
	hyperkit/src/pci_virtio_net_tap.c \
	hyperkit/src/pci_virtio_net_vmnet.c \
    hyperkit/src/pci_virtio_net_vpnkit.c \
	hyperkit/src/pci_virtio_rnd.c \
	hyperkit/src/pm.c \
	hyperkit/src/post.c \
	hyperkit/src/rtc.c \
	hyperkit/src/smbiostbl.c \
	hyperkit/src/task_switch.c \
	hyperkit/src/uart_emul.c \
	hyperkit/src/xhyve.c \
	hyperkit/src/virtio.c \
	hyperkit/src/xmsr.c hyperkit/src/mirage_block_c.h

FIRMWARE_SRC := \
	hyperkit/src/firmware/bootrom.c \
	hyperkit/src/firmware/kexec.c \
	hyperkit/src/firmware/fbsd.c

HAVE_OCAML_QCOW := $(shell if ocamlfind query qcow uri >/dev/null 2>/dev/null ; then echo YES ; else echo NO; fi)

ifeq ($(HAVE_OCAML_QCOW),YES)
CFLAGS += -DHAVE_OCAML=1 -DHAVE_OCAML_QCOW=1 -DHAVE_OCAML=1

OCAML_SRC := \
	hyperkit/src/mirage_block_ocaml.ml

OCAML_C_SRC := \
	hyperkit/src/mirage_block_c.c

OCAML_WHERE := $(shell ocamlc -where)
OCAML_PACKS := cstruct cstruct.lwt io-page io-page.unix uri mirage-block mirage-block-unix qcow unix threads lwt lwt.unix
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

build/xhyve.o: CFLAGS += -I$(OCAML_WHERE)
endif

SRC := \
	$(VMM_SRC) \
	$(XHYVE_SRC) \
	$(FIRMWARE_SRC) \
	$(OCAML_C_SRC)
