package plugin

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/docker/machine/libmachine"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/drivers/rpc"
)

func init() {
	gob.Register(new(rpcdriver.RpcFlags))
}

func RegisterDriver(d drivers.Driver) {
	if len(os.Args) != 1 {
		fmt.Println("Improper number of arguments.  Usage: ./docker-machine-[driver]")
		os.Exit(1)
	}

	libmachine.SetDebug(true)

	rpcd := new(rpcdriver.RpcServerDriver)
	rpcd.ActualDriver = d
	rpcd.CloseCh = make(chan bool)
	rpc.Register(rpcd)

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println(listener.Addr())

	go http.Serve(listener, nil)

	<-rpcd.CloseCh
}
