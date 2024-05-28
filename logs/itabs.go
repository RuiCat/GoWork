package logs

import (
	"reflect"
	"unsafe"
)

//go:linkname iteraiite_itabs runtime.iterate_itabs
func iteraiite_itabs(fn func(unsafe.Pointer))

// IterateItabs 获取系统导出表
func IterateItabs(fn func(reflect.Type)) {
	iteraiite_itabs(func(p unsafe.Pointer) {
		fn(reflect.TypeOf(*(*any)(p)))
	})
}
