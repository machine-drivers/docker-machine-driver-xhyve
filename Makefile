default: build

clean:
	$(RM) ./bin/docker-machine-xhyve
	$(RM) $(GOPATH)/bin/docker-machine-xhyve

build: clean
	GOGC=off go build -i -o ./bin/docker-machine-xhyve ./bin
	sudo chown root:admin ./bin/docker-machine-xhyve
	sudo chmod u+s ./bin/docker-machine-xhyve

install: build
	cp ./bin/docker-machine-xhyve $(GOPATH)/bin/
	sudo chown root:admin $(GOPATH)/bin/docker-machine-xhyve
	sudo chmod u+s $(GOPATH)/bin/docker-machine-xhyve

.PHONY: build install
