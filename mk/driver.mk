DOCKER_MACHINE_CMD := docker-machine
DOCKER_MACHINE_VM_NAME := xhyve-test
DOCKER_MACHINE_VM_DISKSIZE := 2000
DOCKER_MACHINE_STORAGEPATH := $(HOME)/.docker/machine-test

# Always enable debug mode
export MACHINE_DEBUG=1

default: build

driver-test: test-driver-ls

test-driver-env:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} env ${DOCKER_MACHINE_VM_NAME}

test-driver-ip:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} ip ${DOCKER_MACHINE_VM_NAME}

test-driver-kill:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} kill ${DOCKER_MACHINE_VM_NAME}

test-driver-ls:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} ls

test-driver-regenerate-certs:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} regenerate-certs ${DOCKER_MACHINE_VM_NAME}

test-driver-restart:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} restart ${DOCKER_MACHINE_VM_NAME}

test-driver-rm:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} rm ${DOCKER_MACHINE_VM_NAME}

test-driver-ssh:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} ssh ${DOCKER_MACHINE_VM_NAME} 

test-driver-status:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} status ${DOCKER_MACHINE_VM_NAME}

test-driver-start:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} start ${DOCKER_MACHINE_VM_NAME}

test-driver-upgrade:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} upgrade ${DOCKER_MACHINE_VM_NAME}

test-driver-url:
	${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} url ${DOCKER_MACHINE_VM_NAME}

driver-run: clean build install
	rm -rf ${DOCKER_MACHINE_STORAGEPATH}/machines/${DOCKER_MACHINE_VM_NAME} && ${DOCKER_MACHINE_CMD} --storage-path ${DOCKER_MACHINE_STORAGEPATH} create --driver xhyve --xhyve-disk-size ${DOCKER_MACHINE_VM_DISKSIZE} ${DOCKER_MACHINE_VM_NAME}

driver-kill:
	$(shell for i in $$(ps ax | grep goxhyve | awk '{print $$1'}); do sudo kill 2>/dev/null $$i; done)

.PHONY: driver-kill
