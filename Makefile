default: build

clean:
	$(RM) docker-machine-xhyve

build: clean
	GOGC=off go build -i -o docker-machine-xhyve ./bin

install: build
	cp ./docker-machine-xhyve $(GOPATH)/bin/

.PHONY: build install
