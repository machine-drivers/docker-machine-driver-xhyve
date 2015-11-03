default: build

CGO_CFLAGS="-DGUESTFS_PRIVATE=1 -I/usr/local/include"
CGO_LDFLAGS="-L/usr/local/lib -lguestfs"

clean:
	$(RM) ./bin/docker-machine-driver-xhyve
	$(RM) $(GOPATH)/bin/docker-machine-driver-xhyve

bin/docker-machine-driver-xhyve:
	godep go build -i -o ./bin/docker-machine-driver-xhyve ./bin
	sudo chown root:wheel ./bin/docker-machine-driver-xhyve
	sudo chmod u+s ./bin/docker-machine-driver-xhyve

build: bin/docker-machine-driver-xhyve

install: bin/docker-machine-driver-xhyve
	sudo cp -p ./bin/docker-machine-driver-xhyve $(GOPATH)/bin/

.PHONY: clean
