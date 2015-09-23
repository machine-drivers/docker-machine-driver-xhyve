package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/zchee/docker-machine-xhyve"
)

func main() {
	plugin.RegisterDriver(new(xhyve.Driver))
}
