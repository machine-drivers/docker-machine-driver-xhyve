default: build

build: 
	GOGC=off go build -i -o docker-machine-xhyve ./bin

clean:
	$(RM) docker-machine-xhyve

install: build
	cp ./docker-machine-xhyve /usr/local/bin/docker-machine-xhyve

.PHONY: build install
