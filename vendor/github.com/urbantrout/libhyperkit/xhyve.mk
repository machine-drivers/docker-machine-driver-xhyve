XHYVE_VERSION := $(shell cd hyperkit; git describe --abbrev=6 --dirty --always --tags)
GIT_VERSION_SHA1 := $(shell cd hyperkit; git rev-parse HEAD)

VMM_LIB_SRC := \
	hyperkit/src/lib/vmm/intel/vmcs.c \
	hyperkit/src/lib/vmm/intel/vmx.c \
	hyperkit/src/lib/vmm/intel/vmx_msr.c \
	\
	hyperkit/src/lib/vmm/io/vatpic.c \
	hyperkit/src/lib/vmm/io/vatpit.c \
	hyperkit/src/lib/vmm/io/vhpet.c \
	hyperkit/src/lib/vmm/io/vioapic.c \
	hyperkit/src/lib/vmm/io/vlapic.c \
	hyperkit/src/lib/vmm/io/vpmtmr.c \
	hyperkit/src/lib/vmm/io/vrtc.c \
	\
	hyperkit/src/lib/vmm/vmm.c \
	hyperkit/src/lib/vmm/vmm_api.c \
	hyperkit/src/lib/vmm/vmm_callout.c \
	hyperkit/src/lib/vmm/vmm_host.c \
	hyperkit/src/lib/vmm/vmm_instruction_emul.c \
	hyperkit/src/lib/vmm/vmm_ioport.c \
	hyperkit/src/lib/vmm/vmm_lapic.c \
	hyperkit/src/lib/vmm/vmm_mem.c \
	hyperkit/src/lib/vmm/vmm_stat.c \
	hyperkit/src/lib/vmm/vmm_util.c \
	hyperkit/src/lib/vmm/x86.c

HYPERKIT_LIB_SRC := \
	hyperkit/src/lib/acpitbl.c \
	hyperkit/src/lib/atkbdc.c \
	hyperkit/src/lib/block_if.c \
	hyperkit/src/lib/consport.c \
	hyperkit/src/lib/dbgport.c \
	hyperkit/src/lib/inout.c \
	hyperkit/src/lib/ioapic.c \
	hyperkit/src/lib/md5c.c \
	hyperkit/src/lib/mem.c \
	hyperkit/src/lib/mevent.c \
	hyperkit/src/lib/mptbl.c \
	hyperkit/src/lib/pci_ahci.c \
	hyperkit/src/lib/pci_emul.c \
	hyperkit/src/lib/pci_hostbridge.c \
	hyperkit/src/lib/pci_irq.c \
	hyperkit/src/lib/pci_lpc.c \
	hyperkit/src/lib/pci_uart.c \
	hyperkit/src/lib/pci_virtio_9p.c \
	hyperkit/src/lib/pci_virtio_block.c \
	hyperkit/src/lib/pci_virtio_net_tap.c \
	hyperkit/src/lib/pci_virtio_net_vmnet.c \
	hyperkit/src/lib/pci_virtio_net_vpnkit.c \
	hyperkit/src/lib/pci_virtio_rnd.c \
	hyperkit/src/lib/pci_virtio_sock.c \
	hyperkit/src/lib/pm.c \
	hyperkit/src/lib/post.c \
	hyperkit/src/lib/rtc.c \
	hyperkit/src/lib/smbiostbl.c \
	hyperkit/src/lib/task_switch.c \
	hyperkit/src/lib/uart_emul.c \
	hyperkit/src/lib/virtio.c \
	hyperkit/src/lib/xmsr.c

FIRMWARE_LIB_SRC := \
	hyperkit/src/lib/firmware/bootrom.c \
	hyperkit/src/lib/firmware/kexec.c \
	hyperkit/src/lib/firmware/fbsd.c

HYPERKIT_SRC := hyperkit/src/hyperkit.c

HAVE_OCAML_QCOW := $(shell if ocamlfind query qcow uri >/dev/null 2>/dev/null ; then echo YES ; else echo NO; fi)

ifeq ($(HAVE_OCAML_QCOW),YES)
CFLAGS += -DHAVE_OCAML=1 -DHAVE_OCAML_QCOW=1 -DHAVE_OCAML=1

# prefix vsock file names if PRI_ADDR_PREFIX
# is defined. (not applied to aliases)
ifneq ($(PRI_ADDR_PREFIX),)
CFLAGS += -DPRI_ADDR_PREFIX=\"$(PRI_ADDR_PREFIX)\"
endif

# override default connect socket name if 
# CONNECT_SOCKET_NAME is defined 
ifneq ($(CONNECT_SOCKET_NAME),)
CFLAGS += -DCONNECT_SOCKET_NAME=\"$(CONNECT_SOCKET_NAME)\"
endif

OCAML_SRC := \
	hyperkit/src/lib/mirage_block_ocaml.ml

OCAML_C_SRC := \
	hyperkit/src/lib/mirage_block_c.c

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

build/hyperkit.o: CFLAGS += -I$(OCAML_WHERE)
endif

SRC := \
	$(VMM_LIB_SRC) \
	$(HYPERKIT_LIB_SRC) \
	$(FIRMWARE_LIB_SRC) \
	$(OCAML_C_SRC) \
	$(HYPERKIT_SRC)
