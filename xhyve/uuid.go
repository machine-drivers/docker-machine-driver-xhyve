// Copyright 2015 The docker-machine-driver-xhyve Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xhyve

/*
#include <stdio.h>
#include <uuid/uuid.h>

char uuid_str[37];

extern inline char* uuidgen() {
	uuid_t uuid;

	// generate with random
	uuid_generate_random(uuid);

	// unparse (to string)
	uuid_unparse_upper(uuid, uuid_str);

	return uuid_str;
}
*/
import "C"

func uuidgen() string {
	return C.GoString(C.uuidgen())
}
