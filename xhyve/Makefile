.PHONY: clean all

all:
	./generate.sh
	go install

clean:
	@rm -f *.c
	@git apply -R upstream.patch 2>/dev/null || true
