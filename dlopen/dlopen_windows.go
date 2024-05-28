package dlopen

import (
	"syscall"
	_ "unsafe" // only for go:linkname
)

type syscall15Args struct {
	fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15 uintptr
	f1, f2, f3, f4, f5, f6, f7, f8                                       uintptr
	r1, r2, err                                                          uintptr
}

//go:nosplit
func Syscall_syscall15X(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15 uintptr) (r1, r2, err uintptr) {
	r1, r2, errno := syscall.Syscall15(fn, 15, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15)
	return r1, r2, uintptr(errno)
}
