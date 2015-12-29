MAKEFLAGS := -j 1

DOCKER_MACHINE_CMD := docker-machine

DOCKER_MACHINE_STORAGEPATH := $(HOME)/.docker/machine-test
DOCKER_MACHINE_VM_NAME := xhyve-test
# Set boot2docker-url to already exists iso path
DOCKER_MACHINE_VM_BOOT2DOCKER_URL := $(HOME)/.docker/machine/cache/boot2docker.iso
# Set CPU size to hw.ncpu/2
DOCKER_MACHINE_VM_CPU_COUNT := ${shell /usr/bin/python -c "print(${shell sysctl -n hw.ncpu}/2)"}
# Set memory size to hw.memsize/2 MB
DOCKER_MACHINE_VM_MEMORY_SIZE := ${shell /usr/bin/python -c "print(${shell sysctl -n hw.memsize}/2097152)"}
DOCKER_MACHINE_VM_DISKSIZE := 2000

# Always enable debug mode
export MACHINE_DEBUG=1

default: build

test-env:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} env ${DOCKER_MACHINE_VM_NAME}

test-inspect:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} inspect ${DOCKER_MACHINE_VM_NAME}

test-ip:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} ip ${DOCKER_MACHINE_VM_NAME}

test-kill:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} kill ${DOCKER_MACHINE_VM_NAME}

test-ls:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} ls

test-regenerate-certs:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} regenerate-certs ${DOCKER_MACHINE_VM_NAME}

test-restart:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} restart ${DOCKER_MACHINE_VM_NAME}

test-rm:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} rm -f ${DOCKER_MACHINE_VM_NAME}

test-ssh:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} ssh ${DOCKER_MACHINE_VM_NAME}

test-status:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} status ${DOCKER_MACHINE_VM_NAME}

test-stop:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} stop ${DOCKER_MACHINE_VM_NAME}

test-start:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} start ${DOCKER_MACHINE_VM_NAME}

test-upgrade:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} upgrade ${DOCKER_MACHINE_VM_NAME}

test-url:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} url ${DOCKER_MACHINE_VM_NAME}

driver-run: clean build install driver-remove
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} create --driver xhyve \
		--xhyve-boot2docker-url ${DOCKER_MACHINE_VM_BOOT2DOCKER_URL} \
		--xhyve-cpu-count ${DOCKER_MACHINE_VM_CPU_COUNT} \
		--xhyve-memory-size ${DOCKER_MACHINE_VM_MEMORY_SIZE} \
		--xhyve-disk-size ${DOCKER_MACHINE_VM_DISKSIZE} \
		--xhyve-experimental-nfs-share \
		${DOCKER_MACHINE_VM_NAME}

driver-kill:
	@PID=$$(pgrep goxhyve) || PID=none; \
	echo "${CBLUE}==>${CRESET}${CGREEN}Kill goxhyve test process. PID:$$PID ${CRESET} ..."; \
	sudo kill $$PID 2>/dev/null || true

driver-remove:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} rm -f ${DOCKER_MACHINE_VM_NAME} || true
	$(if $(shell ls $(HOME)/.docker/machine-test), sudo rm -rf ${DOCKER_MACHINE_STORAGEPATH},)

.PHONY: driver-kill
