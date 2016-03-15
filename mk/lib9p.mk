MAKEFLAGS := -j1
LIB9P_DIR := vendor/lib9p

LIB9P_CFLAGS := \
	-Weverything \
	-Wno-padded \
	-Wno-gnu-zero-variadic-macro-arguments \
	-Wno-format-nonliteral \
	-Werror \
	-g \
	-O0

LIB9P_SRCS := \
	pack.c \
	connection.c \
	request.c \
	log.c \
	hashtable.c \
	utils.c \
	sbuf/sbuf.c \
	transport/socket.c \
	backend/fs.c

VENDOR_BUILD_DIR := vendor/build
LIB9P_BUILD_DIR := $(VENDOR_BUILD_DIR)/lib9p
LIB9P_LIB_SRCS := $(addprefix ${LIB9P_DIR}/,$(LIB9P_SRCS))
LIB9P_LIB_OBJS := $(addprefix ${LIB9P_BUILD_DIR}/,$(LIB9P_SRCS:.c=.o))
LIB9P_LIB := ${LIB9P_BUILD_DIR}/lib9p.a
LIB9P_DYLIB := ${LIB9P_BUILD_DIR}/lib9p.dylib

default: build

lib9p: vendor/build/lib9p $(LIB9P_LIB)

vendor/lib9p/%.c:
	${GIT_CMD} submodule update --init

vendor/build/lib9p: vendor/lib9p/%.c
	mkdir -p ${LIB9P_BUILD_DIR} ${LIB9P_BUILD_DIR}/sbuf ${LIB9P_BUILD_DIR}/transport ${LIB9P_BUILD_DIR}/backend

vendor/build/lib9p/%.o: vendor/lib9p/%.c
	$(CC) $(LIB9P_CFLAGS) -c $< -o $@

$(LIB9P_LIB): $(LIB9P_LIB_OBJS)
	$(LIBTOOL) -static $^ -o $@

$(LIB9P_DYLIB): $(LIB9P_LIB_OBJS)
	$(CC) -dynamiclib $^ -o ${LIB9P_BUILD_DIR}/$@

clean-lib9p:
	@${RM} -r ${VENDOR_BUILD_DIR}

.PHONY: clean-lib9p lib9p
.PRECIOUS: vendor/lib9p/%.c
