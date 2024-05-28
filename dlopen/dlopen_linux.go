package dlopen

import (
	"errors"
	"fmt"
	"path/filepath"

	"unsafe"
)

/*
#cgo LDFLAGS: -ldl
#include <stdlib.h>
#include <stdint.h>
#include <dlfcn.h>
#include <errno.h>
#include <assert.h>

typedef struct syscall15Args {
	uintptr_t fn;
	uintptr_t a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15;
	uintptr_t f1, f2, f3, f4, f5, f6, f7, f8;
	uintptr_t r1, r2, err;
} syscall15Args;

void syscall15(struct syscall15Args *args) {
	assert((args->f1|args->f2|args->f3|args->f4|args->f5|args->f6|args->f7|args->f8) == 0);
	uintptr_t (*func_name)(uintptr_t a1, uintptr_t a2, uintptr_t a3, uintptr_t a4, uintptr_t a5, uintptr_t a6,
		uintptr_t a7, uintptr_t a8, uintptr_t a9, uintptr_t a10, uintptr_t a11, uintptr_t a12,
		uintptr_t a13, uintptr_t a14, uintptr_t a15);
	*(void**)(&func_name) = (void*)(args->fn);
	uintptr_t r1 =  func_name(args->a1,args->a2,args->a3,args->a4,args->a5,args->a6,args->a7,args->a8,args->a9,
		args->a10,args->a11,args->a12,args->a13,args->a14,args->a15);
	args->r1 = r1;
	args->err = errno;
}
*/
import "C"

// LibHandle represents an open handle to a library (.so)
type libHandle struct {
	Handle  unsafe.Pointer
	Libname string
}

// GetHandle tries to get a handle to a library (.so), attempting to access it
// by the names specified in libs and returning the first that is successfully
// opened. Callers are responsible for closing the handler. If no library can
// be successfully opened, an error is returned.
func GetHandle(libs ...string) (Handle, error) {
	for _, name := range libs {
		name, _ = filepath.Abs(name)
		libname := C.CString(name)
		defer C.free(unsafe.Pointer(libname))
		handle := C.dlopen(libname, C.RTLD_LAZY|C.RTLD_GLOBAL)
		if handle != nil {
			h := &libHandle{
				Handle:  handle,
				Libname: name,
			}
			return h, nil
		}
	}
	return nil, errors.New(C.GoString(C.dlerror()))
}
func (l *libHandle) GetHandle() unsafe.Pointer {
	return l.Handle
}

// GetSymbolPointer takes a symbol name and returns a pointer to the symbol.
func (l *libHandle) GetSymbolPointer(symbol string) (unsafe.Pointer, error) {
	sym := C.CString(symbol)
	defer C.free(unsafe.Pointer(sym))

	C.dlerror()
	p := C.dlsym(l.Handle, sym)
	e := C.dlerror()
	if e != nil {
		return nil, fmt.Errorf("error resolving symbol %q: %v", symbol, errors.New(C.GoString(e)))
	}

	return p, nil
}

// Close closes a libHandle.
func (l *libHandle) Close() error {
	C.dlerror()
	C.dlclose(l.Handle)
	e := C.dlerror()
	if e != nil {
		return fmt.Errorf("error closing %v: %v", l.Libname, errors.New(C.GoString(e)))
	}
	return nil
}

// this is only here to make the assembly files happy :)
type syscall15Args struct {
	fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15 uintptr
	f1, f2, f3, f4, f5, f6, f7, f8                                       uintptr
	r1, r2, err                                                          uintptr
}

//go:nosplit
func syscall15X(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15 uintptr) (r1, r2, err uintptr) {
	args := C.syscall15Args{
		C.uintptr_t(fn), C.uintptr_t(a1), C.uintptr_t(a2), C.uintptr_t(a3),
		C.uintptr_t(a4), C.uintptr_t(a5), C.uintptr_t(a6),
		C.uintptr_t(a7), C.uintptr_t(a8), C.uintptr_t(a9), C.uintptr_t(a10), C.uintptr_t(a11), C.uintptr_t(a12),
		C.uintptr_t(a13), C.uintptr_t(a14), C.uintptr_t(a15), 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
	C.syscall15(&args)
	return uintptr(args.r1), uintptr(args.r2), uintptr(args.err)
}

//go:nosplit
func Syscall_syscall15X(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15 uintptr) (r1, r2, err uintptr) {
	return syscall15X(fn, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15)
}
