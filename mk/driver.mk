JOBS=1

DOCKER_MACHINE_CMD := docker-machine
DOCKER_MACHINE_VM_NAME := xhyve-test
DOCKER_MACHINE_VM_DISKSIZE := 2000
DOCKER_MACHINE_STORAGEPATH := $(HOME)/.docker/machine-test

# Always enable debug mode
export MACHINE_DEBUG=1

default: build

driver-test: test-ls

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
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} rm ${DOCKER_MACHINE_VM_NAME}

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

driver-run: clean build install driver-kill
	rm -rf ${DOCKER_MACHINE_STORAGEPATH}/machines/${DOCKER_MACHINE_VM_NAME} && ${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} create --driver xhyve --xhyve-disk-size ${DOCKER_MACHINE_VM_DISKSIZE} ${DOCKER_MACHINE_VM_NAME}

driver-kill:
	@PID=$$(pgrep goxhyve) || PID=none; \
	echo "${CBLUE}==>${CRESET}${CGREEN}Kill goxhyve test process. PID:$$PID ${CRESET} ..."; \
	sudo kill $$PID 2>/dev/null || true

.PHONY: driver-kill
