package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/nathanleclaire/docker-machine-xhyve"
)

func main() {
	plugin.RegisterDriver(new(xhyve.Driver))
}
