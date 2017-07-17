package main

import (
	"fmt"
	"os"
	"syscall"

	vmnet "github.com/zchee/go-vmnet"
)

func main() {
	binary, err := os.Stat(os.Args[0])
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		panic(fmt.Sprintf("Usage: %s <UUID>", binary.Name()))
	}

	if syscall.Getuid() != 0 {
		panic(fmt.Sprintf("%s needs sudo privileges. Please add 'sudo' before the command", binary.Name()))
	}

	uuid := os.Args[1]
	mac, err := vmnet.GetMACAddressFromUUID(uuid)
	if err != nil {
		panic(fmt.Sprintf("vmnet: error from vmnet with %+v", err))
	}

	fmt.Println(mac)
}
