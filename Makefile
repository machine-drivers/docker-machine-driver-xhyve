default: build

clean:
	$(RM) docker-machine-xhyve
	$(RM) $(GOPATH)/bin/docker-machine-xhyve

build: clean
	GOGC=off go build -i -o docker-machine-xhyve ./bin
	sudo chown root ./docker-machine-xhyve
	sudo chmod +s ./docker-machine-xhyve 
install: build
	cp ./docker-machine-xhyve $(GOPATH)/bin/

.PHONY: build install
