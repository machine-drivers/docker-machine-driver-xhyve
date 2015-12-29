package main

import (
	"os"

	"github.com/docker/machine/libmachine/drivers/plugin"
	goxhyve "github.com/hooklift/xhyve"
	"github.com/zchee/docker-machine-driver-xhyve"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "xhyve" {
		args := os.Args[1:] // Technically we only need 2:, but a bug in hooklift/xhyve requires one arg in the beginning
		if err := goxhyve.Run(args, make(chan string)); err != nil {
			panic(err)
		}
	} else {
		plugin.RegisterDriver(xhyve.NewDriver("", ""))
	}
}
