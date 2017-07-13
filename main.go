// Copyright 2015 The docker-machine-driver-xhyve Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/zchee/docker-machine-driver-xhyve/xhyve"
)

func main() {
	plugin.RegisterDriver(xhyve.NewDriver("", ""))
}
