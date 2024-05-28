package basebyte

import (
	"encoding/binary"
	"math"
)

// Read 读
type Read struct {
	Byte   *[]byte
	Offset int
	Order  binary.ByteOrder
}

// Bool 逻辑型
func (r *Read) Bool() (v bool) {
	v = (*r.Byte)[r.Offset] != 0
	r.Offset++
	return
}

// Int8 单字节整数
func (r *Read) Int8() (v int8) {
	v = int8((*r.Byte)[r.Offset])
	r.Offset++
	return v
}

// Uint8  单字节正整数
func (r *Read) Uint8() (v uint8) {
	v = (*r.Byte)[r.Offset]
	r.Offset++
	return v
}

// Int16 双字节整数
func (r *Read) Int16() (v int16) {
	v = int16(r.Order.Uint16((*r.Byte)[r.Offset:]))
	r.Offset += 2
	return v
}

// Uint16 双字节正整数
func (r *Read) Uint16() (v uint16) {
	v = r.Order.Uint16((*r.Byte)[r.Offset:])
	r.Offset += 2
	return v
}

// Int32 三字节整数
func (r *Read) Int32() (v int32) {
	v = int32(r.Order.Uint32((*r.Byte)[r.Offset:]))
	r.Offset += 4
	return v
}

// Uint32 三字节正整数
func (r *Read) Uint32() (v uint32) {
	v = r.Order.Uint32((*r.Byte)[r.Offset:])
	r.Offset += 4
	return v
}

// Int64 四字节整数
func (r *Read) Int64() (v int64) {
	v = int64(r.Order.Uint64((*r.Byte)[r.Offset:]))
	r.Offset += 8
	return v
}

// Uint64 四字节正整数
func (r *Read) Uint64() (v uint64) {
	v = r.Order.Uint64((*r.Byte)[r.Offset:])
	r.Offset += 8
	return v
}

// Int 平台相关整数
func (r *Read) Int() (v int) {
	if bit64 {
		v = int(r.Int64())
	} else {
		v = int(r.Int32())
	}
	return v
}

// Uint 平台相关正整数
func (r *Read) Uint() (v uint) {
	if bit64 {
		v = uint(r.Uint64())
	} else {
		v = uint(r.Uint32())
	}
	return v
}

// Float32 浮点数
func (r *Read) Float32() (v float32) {
	return math.Float32frombits(r.Uint32())
}

// Float64 浮点数
func (r *Read) Float64() (v float64) {
	return math.Float64frombits(r.Uint64())
}

// Complex64 复数
func (r *Read) Complex64() (v complex64) {
	return complex(r.Float32(), r.Float32())
}

// Complex128 复数
func (r *Read) Complex128() (v complex128) {
	return complex(r.Float64(), r.Float64())
}

// Bytes 字节集
func (r *Read) Bytes() (v []byte) {
	i := r.Int()
	v = append(v, (*r.Byte)[r.Offset:r.Offset+i]...)
	r.Offset += i
	return v
}

// GetByte 得到字节集
func (r *Read) GetByte(i int) (v []byte) {
	v = append(v, (*r.Byte)[r.Offset:r.Offset+i]...)
	r.Offset += i
	return v
}
