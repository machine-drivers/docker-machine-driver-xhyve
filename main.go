// Copyright 2015 The docker-machine-driver-xhyve Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/zchee/docker-machine-driver-xhyve/xhyve"
	hyperkit "github.com/zchee/libhyperkit"
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
		if err := hyperkit.Run(args, ptyCh); err != nil {
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
