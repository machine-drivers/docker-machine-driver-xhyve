include xhyve.mk

# Downloads xhyve files to make this library Go gettable.
# It also applies a patch to rename the main function so we can use xhyve
# as a library instead.
upstream:
	git clone https://github.com/mist64/xhyve.git vendor/xhyve
	-cd vendor/xhyve; patch -N -p1 < ../../xhyve.patch
	find . \( -name \*.orig -o -name \*.rej \) -delete
	for file in $(SRC); do \
		cp -f $$file $$(basename $$file) ; \
	done
	cp -r vendor/xhyve/include include

clean:
	rm -rf *.c vendor include

.PHONY: upstream clean
