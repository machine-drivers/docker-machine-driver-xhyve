include xhyve.mk

build:
	go build -o xhyve cmd/xhyve/main.go

clone-xhyve:
	-git clone https://github.com/xhyve-xyz/xhyve.git vendor/xhyve

sync: clone-xhyve apply-patch
	find . \( -name \*.orig -o -name \*.rej \) -delete
	for file in $(SRC); do \
		cp -f $$file $$(basename $$file) ; \
	done
	cp -r vendor/xhyve/include include

apply-patch:
	-cd vendor/xhyve; patch -N -p1 < ../../xhyve.patch

generate-patch: apply-patch
	cd vendor/xhyve; git diff > ../../xhyve.patch

clean:
	rm -rf *.c vendor include

.PHONY: build clone-xhyve sync apply-patch generate-patch clean
