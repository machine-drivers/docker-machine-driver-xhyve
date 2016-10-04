package mem

import "unsafe"

//go:noescape
func Set(dst []byte, value byte)

//go:noescape
func move(to uintptr, from uintptr, n uintptr)

// memmove copies n bytes from "from" to "to".
// in memmove_*.s
func Move(dst, src []byte) {
	move(uintptr(unsafe.Pointer(&dst[0])), uintptr(unsafe.Pointer(&src[0])), uintptr(len(src)))
}

//go:noescape
func cpy(to uintptr, from uintptr, n uintptr)

func Cpy(dst, src []byte, n uintptr) {
	if n == 0 {
		return
	}
	cpy(uintptr(unsafe.Pointer(&dst[0])), uintptr(unsafe.Pointer(&src[0])), n)
}
