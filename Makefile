default: build

clean:
	$(RM) ./bin/docker-machine-driver-xhyve
	$(RM) $(GOPATH)/bin/docker-machine-driver-xhyve

bin/docker-machine-driver-xhyve:
	# ./scripts/static.sh godep go build -i -o ./bin/docker-machine-xhyve ./bin
	godep go build -i -o ./bin/docker-machine-driver-xhyve ./bin
	sudo chown root:admin ./bin/docker-machine-driver-xhyve
	sudo chmod u+s ./bin/docker-machine-driver-xhyve

build: bin/docker-machine-driver-xhyve

install: bin/docker-machine-driver-xhyve
	sudo cp ./bin/docker-machine-driver-xhyve $(GOPATH)/bin/
	sudo chown root:admin $(GOPATH)/bin/docker-machine-driver-xhyve
	sudo chmod u+s $(GOPATH)/bin/docker-machine-driver-xhyve

.PHONY: clean
