package basebyte

import (
	"unsafe"
)

var intSize int = 32 << (^uint(0) >> 63) / 8

func getDataSize(in interface{}) int {
	switch in.(type) {
	case bool, int8, uint8:
		return 1
	case int:
		return intSize
	case int16, uint16:
		return 2
	case int32, uint32, float32:
		return 4
	case int64, uint64, float64, complex64:
		return 8
	case complex128:
		return 16
	case *bool, *int8, *uint8:
		return 1
	case *int:
		return intSize
	case *int16, *uint16:
		return 2
	case *int32, *uint32, *float32:
		return 4
	case *int64, *uint64, *float64, *complex64:
		return 8
	case *complex128:
		return 16
	}
	return -1
}

func getByte(size int, pointer uintptr) []byte {
	b := make([]byte, size)
	for i := 0; i < size; i++ {
		b[i] = *(*byte)(unsafe.Pointer(pointer + uintptr(i)))
	}
	return b
}
func setByte(size int, pointer uintptr, b []byte) int {
	for i := 0; i < size; i++ {
		*(*byte)(unsafe.Pointer(pointer + uintptr(i))) = b[i]
	}
	return size
}

// GetType 序列化
func GetType(in interface{}) []byte {
	size := getDataSize(in)
	if size > 0 {
		pointer := *(*[2]uintptr)(unsafe.Pointer(&in))
		return getByte(size, pointer[1])
	}
	switch v := in.(type) {
	case string:
		return []byte(v)
	case []byte:
		return v
	default:
		return nil
	}
}

// SetType 反序列化
func SetType(in interface{}, b []byte) int {
	size := getDataSize(in)
	if size > 0 {
		pointer := *(*[2]uintptr)(unsafe.Pointer(&in))
		return setByte(size, pointer[1], b)
	}
	switch v := in.(type) {
	case *string:
		*v = string(b)
		return len(b)
	case *[]byte:
		l := len(*v)
		for i := 0; i < l; i++ {
			(*v)[i] = b[i]
		}
		return l
	default:
		return -1
	}
}
