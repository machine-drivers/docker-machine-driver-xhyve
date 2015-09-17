package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tiborvass/xhyve-bindings"
)

func usage() {
	fmt.Println("Usage: ./main <vmlinuz> <initrd>")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 3 {
		usage()
	}
	vmlinuz := os.Args[1]
	initrd := os.Args[2]
	args := strings.Fields("-A -s 0:0,hostbridge -s 31,lpc -l com1 -s 2:0,virtio-net -m 1024M")
	if err := xhyve.Exec(append(args,
		"-U", fmt.Sprintf("%s", uuid),
		fmt.Sprintf("-s 3,ahci-cd,%s", iso),
		fmt.Sprintf("-s 4,virtio-blk,%s", img),
		fmt.Sprintf("-s 5,virtio-blk,%s", userdata),
		"-f", fmt.Sprintf("kexec,%s,%s,loglevel=3 user=docker console=ttyS0 console=tty0 noembed nomodeset norestore waitusb=10:LABEL=boot2docker-data base host=boot2docker", vmlinuz, initrd))...); err != nil {
		log.Fatal(err)
	}
}
