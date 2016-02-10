MAKEFLAGS := -j 1

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

# build/pack.o build/connection.o build/request.o build/log.o build/hashtable.o build/utils.o build/sbuf/sbuf.o build/transport/socket.o build/backend/fs.o -o build/lib9p.dylib

LIB9P_BUILD_DIR := ${LIB9P_DIR}/build
LIB9P_LIB_SRCS := $(addprefix ${LIB9P_DIR}/,$(LIB9P_SRCS))
LIB9P_LIB_OBJS := $(addprefix ${LIB9P_BUILD_DIR}/,$(LIB9P_SRCS:.c=.o))
LIB9P_LIB := lib9p.a
LIB9P_DYLIB := lib9p.dylib

lib9p: $(LIB9P_LIB)
	
$(LIB9P_LIB): $(LIB9P_LIB_OBJS)
	libtool -static $^ -o ${LIB9P_BUILD_DIR}/$@

$(LIB9P_DYLIB): $(LIB9P_LIB_OBJS)
	cc -dynamiclib $^ -o ${LIB9P_BUILD_DIR}/$@
	
clean-lib9p:
	rm -rf ${LIB9P_DIR}

${LIB9P_LIB_SRCS}:
	git submodule update --init
	mkdir -p ${LIB9P_BUILD_DIR} ${LIB9P_BUILD_DIR}/sbuf ${LIB9P_BUILD_DIR}/transport ${LIB9P_BUILD_DIR}/backend

vendor/lib9p/build/%.o: vendor/lib9p/%.c
	$(CC) $(LIB9P_CFLAGS) -c $< -o $@
