// Package xhyve allows to embed xhyve in Go.
package xhyve

import (
	"fmt"
	"runtime"
	"unsafe"
)

/*
#cgo CFLAGS: -I upstream/include/ -DMAIN=xhyve_main
#cgo LDFLAGS: -framework Hypervisor -framework vmnet

#include "xhyve/xhyve.h"
#include <stdlib.h>
*/
import "C"

// Exec will call xhyve's main function passing `xhyve`
// as `argv[0]` and `args` for the rest.
//
// Example: Exec("-v") will set argc to 2, and argv to ["xhyve", "-v"].
func Exec(args ...string) error {
	args = append([]string{"xhyve"}, args...)

	var x *C.char
	size := int(unsafe.Sizeof(x))
	argv := C.malloc(C.size_t(len(args) * size))
	defer C.free(unsafe.Pointer(argv))
	for i, arg := range args {
		ptr := unsafe.Pointer(uintptr(argv) + uintptr(size*i))
		*(**C.char)(ptr) = C.CString(arg)
	}

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if e := C.xhyve_main(C.int(len(args)), (**C.char)(argv)); e != 0 {
		return fmt.Errorf("xhyve error: %d", e)
	}

	return nil
}
