package dlopen

import (
	"math"
	"reflect"
	"runtime"
	"unsafe"
)

/*
#cgo LDFLAGS: -ldl
#include <stdlib.h>
#include <stdint.h>
*/
import "C"

const (
	maxArgs     = 15
	numOfFloats = 8 // arm64 and amd64 both have 8 float registers
)

//go:linkname runtime_noescape runtime.noescape
//go:noescape
func runtime_noescape(p unsafe.Pointer) unsafe.Pointer // from runtime/stubs.go

// Handle 动态库接口
type Handle interface {
	GetHandle() unsafe.Pointer
	GetSymbolPointer(symbol string) (unsafe.Pointer, error)
	Close() error
}

// GetValue 得到值
func GetValue[T any](handle Handle, symbol string) *T {
	ptr, err := handle.GetSymbolPointer(symbol)
	if err != nil {
		panic(err)
	}
	return (*T)(ptr)
}

// GetFn 得到调用函数
func GetFn[Fn any](fnPtr unsafe.Pointer, typePtr uintptr) Fn {
	ty := reflect.TypeOf((*Fn)(nil)).Elem()
	if ty.Kind() != reflect.Func {
		panic("purego: fptr must be a function pointer")
	}
	if ty.NumOut() > 1 {
		panic("purego: function can only return zero or one values")
	}
	if fnPtr == nil {
		panic("purego: cfn is nil")
	}
	{
		// this code checks how many registers and stack this function will use
		// to avoid crashing with too many arguments
		var ints int
		var floats int
		var stack int
		if typePtr != 0 {
			if ints < numOfIntegerRegisters() {
				ints++
			} else {
				stack++
			}
		}
		for i := 0; i < ty.NumIn(); i++ {
			arg := ty.In(i)
			switch arg.Kind() {
			case reflect.String, reflect.Uintptr, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Ptr, reflect.UnsafePointer, reflect.Slice,
				reflect.Func, reflect.Bool:
				if ints < numOfIntegerRegisters() {
					ints++
				} else {
					stack++
				}
			case reflect.Float32, reflect.Float64:
				if floats < numOfFloats {
					floats++
				} else {
					stack++
				}
			default:
				panic("purego: unsupported kind " + arg.Kind().String())
			}
		}
		sizeOfStack := maxArgs - numOfIntegerRegisters()
		if stack > sizeOfStack {
			panic("purego: too many arguments")
		}
	}
	return (reflect.MakeFunc(ty, func(args []reflect.Value) (results []reflect.Value) {
		if len(args) > 0 {
			if variadic, ok := args[len(args)-1].Interface().([]interface{}); ok {
				// subtract one from args bc the last argument in args is []interface{}
				// which we are currently expanding
				tmp := make([]reflect.Value, len(args)-1+len(variadic))
				n := copy(tmp, args[:len(args)-1])
				for i, v := range variadic {
					tmp[n+i] = reflect.ValueOf(v)
				}
				args = tmp
			}
		}
		var sysargs [maxArgs]uintptr
		stack := sysargs[numOfIntegerRegisters():]
		var floats [numOfFloats]uintptr
		var numInts int
		var numFloats int
		var numStack int
		var addStack, addInt, addFloat func(x uintptr)
		if runtime.GOARCH == "arm64" || runtime.GOOS != "windows" {
			// Windows arm64 uses the same calling convention as macOS and Linux
			addStack = func(x uintptr) {
				stack[numStack] = x
				numStack++
			}
			addInt = func(x uintptr) {
				if numInts >= numOfIntegerRegisters() {
					addStack(x)
				} else {
					sysargs[numInts] = x
					numInts++
				}
			}
			addFloat = func(x uintptr) {
				if numFloats < len(floats) {
					floats[numFloats] = x
					numFloats++
				} else {
					addStack(x)
				}
			}
		} else {
			// On Windows amd64 the arguments are passed in the numbered registered.
			// So the first int is in the first integer register and the first float
			// is in the second floating register if there is already a first int.
			// This is in contrast to how macOS and Linux pass arguments which
			// tries to use as many registers as possible in the calling convention.
			addStack = func(x uintptr) {
				sysargs[numStack] = x
				numStack++
			}
			addInt = addStack
			addFloat = addStack
		}
		var keepAlive []interface{}
		defer func() {
			runtime.KeepAlive(keepAlive)
			runtime.KeepAlive(args)
		}()
		if typePtr != 0 {
			addInt(typePtr)
		}
		for _, v := range args {
			switch v.Kind() {
			case reflect.String:
				ptr := CString(v.String())
				keepAlive = append(keepAlive, ptr)
				addInt(uintptr(unsafe.Pointer(ptr)))
			case reflect.Uintptr, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				addInt(uintptr(v.Uint()))
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				addInt(uintptr(v.Int()))
			case reflect.Ptr, reflect.UnsafePointer, reflect.Slice:
				// There is no need to keepAlive this pointer separately because it is kept alive in the args variable
				addInt(v.Pointer())
			case reflect.Bool:
				if v.Bool() {
					addInt(1)
				} else {
					addInt(0)
				}
			case reflect.Float32:
				addFloat(uintptr(math.Float32bits(float32(v.Float()))))
			case reflect.Float64:
				addFloat(uintptr(math.Float64bits(v.Float())))
			default:
				panic("purego: unsupported kind: " + v.Kind().String())
			}
		}

		// This is a fallback for Windows amd64, 386, and arm. Note this may not support floats
		r1, r2, _ := Syscall_syscall15X(uintptr(fnPtr), sysargs[0], sysargs[1], sysargs[2], sysargs[3], sysargs[4],
			sysargs[5], sysargs[6], sysargs[7], sysargs[8], sysargs[9], sysargs[10], sysargs[11],
			sysargs[12], sysargs[13], sysargs[14])

		if ty.NumOut() == 0 {
			return nil
		}
		outType := ty.Out(0)
		v := reflect.New(outType).Elem()
		switch outType.Kind() {
		case reflect.Uintptr, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v.SetUint(uint64(r1))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(int64(r1))
		case reflect.Bool:
			v.SetBool(byte(r1) != 0)
		case reflect.UnsafePointer:
			// We take the address and then dereference it to trick go vet from creating a possible miss-use of unsafe.Pointer
			v.SetPointer(*(*unsafe.Pointer)(unsafe.Pointer(&r1)))
		case reflect.Ptr:
			// It is safe to have the address of r1 not escape because it is immediately dereferenced with .Elem()
			v = reflect.NewAt(outType, runtime_noescape(unsafe.Pointer(&r1))).Elem()
		case reflect.String:
			v.SetString(GoString(r1))
		case reflect.Float32:
			// NOTE: r2 is only the floating return value on 64bit platforms.
			// On 32bit platforms r2 is the upper part of a 64bit return.
			v.SetFloat(float64(math.Float32frombits(uint32(r2))))
		case reflect.Float64:
			// NOTE: r2 is only the floating return value on 64bit platforms.
			// On 32bit platforms r2 is the upper part of a 64bit return.
			v.SetFloat(math.Float64frombits(uint64(r2)))
		default:
			panic("purego: unsupported return kind: " + outType.Kind().String())
		}
		return []reflect.Value{v}
	}).Interface()).(Fn)
}

// GetSymbolFn 得到动态库函数
func GetSymbolFn[Fn any](handle Handle, symbol string) Fn {
	return GetSymbolTypeFn[Fn](handle, nil, symbol)
}

// GetSymbolTypeFn 得到动态库结构体函数
func GetSymbolTypeFn[Fn any](handle Handle, typePtr unsafe.Pointer, symbol string) Fn {
	handlePtr, err := handle.GetSymbolPointer(symbol)
	if err != nil {
		panic(err)
	}
	return GetFn[Fn](handlePtr, uintptr(typePtr))
}

func numOfIntegerRegisters() int {
	switch runtime.GOARCH {
	case "arm64":
		return 8
	case "amd64":
		return 6
	// TODO: figure out why 386 tests are not working
	/*case "386":
		return 0
	case "arm":
		return 4*/
	default:
		panic("purego: unknown GOARCH (" + runtime.GOARCH + ")")
	}
}

// hasSuffix tests whether the string s ends with suffix.
func hasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// CString converts a go string to *byte that can be passed to C code.
func CString(name string) *byte {
	if hasSuffix(name, "\x00") {
		return &(*(*[]byte)(unsafe.Pointer(&name)))[0]
	}
	b := make([]byte, len(name)+1)
	copy(b, name)
	return &b[0]
}

// GoString copies a null-terminated char* to a Go string.
func GoString(c uintptr) string {
	// We take the address and then dereference it to trick go vet from creating a possible misuse of unsafe.Pointer
	ptr := *(*unsafe.Pointer)(unsafe.Pointer(&c))
	if ptr == nil {
		return ""
	}
	var length int
	for {
		if *(*byte)(unsafe.Add(ptr, uintptr(length))) == '\x00' {
			break
		}
		length++
	}
	return string(unsafe.Slice((*byte)(ptr), length))
}
