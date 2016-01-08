package main

import (
	"fmt"
	"os"

	"github.com/docker/machine/libmachine/drivers/plugin"
	goxhyve "github.com/hooklift/xhyve"
	"github.com/zchee/docker-machine-driver-xhyve/xhyve"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "xhyve" {
		runXhyve()
	} else {
		plugin.RegisterDriver(xhyve.NewDriver("", ""))
	}
}

func runXhyve() {
	done := make(chan bool)
	ptyCh := make(chan string)

	args := os.Args[1:] // Technically we only need 2:, but a bug in hooklift/xhyve requires one arg in the beginning
	go func() {
		if err := goxhyve.Run(args, ptyCh); err != nil {
			fmt.Println(err)
		}
		done <- true
	}()

	if os.Args[len(os.Args)-1] != "-M" {
		fmt.Printf("Waiting on a pseudo-terminal to be ready... ")
		pty := <-ptyCh
		fmt.Printf("done\n")
		fmt.Printf("Hook up your terminal emulator to %s in order to connect to your VM\n", pty)
	}

	<-done
}
