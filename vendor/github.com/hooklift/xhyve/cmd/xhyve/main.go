package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/hooklift/xhyve"
)

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		s := <-c
		fmt.Printf("Signal %s received\n", s)
	}()
}

func main() {
	done := make(chan bool)
	ptyCh := make(chan string)

	go func() {
		if err := xhyve.Run(os.Args, ptyCh); err != nil {
			fmt.Println(err)
		}
		done <- true
	}()

	fmt.Printf("Waiting on a pseudo-terminal to be ready... ")
	pty := <-ptyCh
	fmt.Printf("done\n")
	fmt.Printf("Hook up your terminal emulator to %s in order to connect to your VM\n", pty)

	<-done
	fmt.Println("Hypervisor goroutine finished!")
}
