include xhyve.mk

build:
	go build -o blob cmd/xhyve/main.go

clone-xhyve:
	-git clone https://github.com/docker/hyperkit.git hyperkit
	# cherry-picked from https://github.com/mist64/xhyve/pull/81
	# Fix non-deterministic delays when accessing a vcpu in "running" or "sleeping" state.
	-cd hyperkit; curl -Ls https://patch-diff.githubusercontent.com/raw/mist64/xhyve/pull/81.patch | patch -N -p1
	# experimental support for raw devices - https://github.com/mist64/xhyve/pull/80
	-cd hyperkit; curl -Ls https://patch-diff.githubusercontent.com/raw/mist64/xhyve/pull/80.patch | patch -N -p1

sync: clean clone-xhyve apply-patch
	find . \( -name \*.orig -o -name \*.rej \) -delete
	for file in $(SRC); do \
		cp -f $$file $$(basename $$file) ; \
		rm -rf $$file ; \
	done
	cp -r hyperkit/include include
	cp hyperkit/README.md README.hyperkit.md
	cp hyperkit/README.xhyve.md .

apply-patch:
	-cd hyperkit; patch -Nl -p1 -F4 < ../xhyve.patch

generate-patch: apply-patch
	cd hyperkit; git diff > ../xhyve.patch

clean:
	rm -rf *.c hyperkit blob include

.PHONY: build clone-xhyve sync apply-patch generate-patch clean
