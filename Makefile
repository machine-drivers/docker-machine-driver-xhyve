default: build

clean:
	$(RM) docker-machine-xhyve
	$(RM) $(GOPATH)/bin/docker-machine-xhyve

build: clean
	GOGC=off go build -i -o docker-machine-xhyve ./bin

install: build
	cp ./docker-machine-xhyve $(GOPATH)/bin/
	sudo chown root $(GOPATH)/bin/docker-machine-xhyve
	sudo chmod +s $(GOPATH)/bin/docker-machine-xhyve

.PHONY: build install
