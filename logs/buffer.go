package logs

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"sync"
	"time"
)

const _size = 1024 // 默认缓冲区大小

// Buffer 数据缓冲
type Buffer struct {
	*bytes.Buffer
	pool Pool
}

// Free 释放
func (b *Buffer) Free() {
	b.Reset()     // 重置
	b.pool.put(b) // 释放
}

// WriteType 写入基础类型
func (b *Buffer) WriteTypeString(a any) {
	switch v := a.(type) {
	case error:
		b.WriteString(v.Error())
	case *error:
		b.WriteTypeString(*v)
	case string:
		b.WriteString(v)
	case *string:
		b.WriteTypeString(*v)
	case bool:
		if v {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	case *bool:
		b.WriteTypeString(*v)
	case float64:
		switch {
		case math.IsNaN(v):
			b.WriteString(`"NaN"`)
		case math.IsInf(v, 1):
			b.WriteString(`"+Inf"`)
		case math.IsInf(v, -1):
			b.WriteString(`"-Inf"`)
		default:
			b.Write(strconv.AppendFloat([]byte{}, v, 'f', -1, 64))
		}
	case *float64:
		b.WriteTypeString(*v)
	case float32:
		switch {
		case math.IsNaN(float64(v)):
			b.WriteString(`"NaN"`)
		case math.IsInf(float64(v), 1):
			b.WriteString(`"+Inf"`)
		case math.IsInf(float64(v), -1):
			b.WriteString(`"-Inf"`)
		default:
			b.Write(strconv.AppendFloat([]byte{}, float64(v), 'f', -1, 32))
		}
	case *float32:
		b.WriteTypeString(*v)
	case int:
		b.WriteTypeString(int64(v))
	case int64:
		b.Write(strconv.AppendInt([]byte{}, v, 10))
	case int32:
		b.WriteTypeString(int64(v))
	case int16:
		b.WriteTypeString(int64(v))
	case int8:
		b.WriteTypeString(int64(v))
	case uint:
		b.WriteTypeString(uint64(v))
	case uint64:
		b.Write(strconv.AppendUint([]byte{}, v, 10))
	case uint32:
		b.WriteTypeString(uint64(v))
	case uint16:
		b.WriteTypeString(uint64(v))
	case uint8:
		b.WriteTypeString(uint64(v))
	case *int:
		b.WriteTypeString(int64(*v))
	case *int64:
		b.WriteTypeString(int64(*v))
	case *int32:
		b.WriteTypeString(int64(*v))
	case *int16:
		b.WriteTypeString(int64(*v))
	case *int8:
		b.WriteTypeString(int64(*v))
	case *uint:
		b.WriteTypeString(uint64(*v))
	case *uint64:
		b.WriteTypeString(uint64(*v))
	case *uint32:
		b.WriteTypeString(uint64(*v))
	case *uint16:
		b.WriteTypeString(uint64(*v))
	case *uint8:
		b.WriteTypeString(uint64(*v))
	case time.Time:
		b.WriteTypeString(v.UnixNano())
	case *time.Time:
		b.WriteTypeString(v.UnixNano())
	case complex128:
		r, i := float64(real(v)), float64(imag(v))
		b.WriteTypeString(r)
		b.WriteByte('+')
		b.WriteTypeString(i)
		b.WriteByte('i')
	case *complex128:
		b.WriteTypeString(*v)
	case complex64:
		b.WriteTypeString(complex128(v))
	case *complex64:
		b.WriteTypeString(complex128(*v))
	case uintptr:
		b.WriteTypeString(uint64(v))
	case *uintptr:
		b.WriteTypeString(uint64(*v))
	default:
		if v, ok := v.(interface{ String() string }); ok {
			b.WriteString(v.String())
		} else {
			b.WriteString("[Type:")
			b.WriteString(reflect.TypeOf(a).String())
			b.WriteString(" Value:")
			b.WriteString(fmt.Sprintf("%v", a))
			b.WriteString("]")
		}
	}
}

// WriteType 写入基础类型
func (b *Buffer) WriteType(a any) bool {
	switch v := a.(type) {
	case *string:
		binary.Write(b, binary.BigEndian, int64(len(*v)))
		b.WriteString(*v)
	case string:
		binary.Write(b, binary.BigEndian, int64(len(v)))
		b.WriteString(v)
	case []string:
		binary.Write(b, binary.BigEndian, int64(len(v)))
		for _, str := range v {
			binary.Write(b, binary.BigEndian, int64(len(str)))
			b.WriteString(str)
		}
	case *int:
		binary.Write(b, binary.BigEndian, int64(*v))
	case int:
		binary.Write(b, binary.BigEndian, int64(v))
	case *uint:
		binary.Write(b, binary.BigEndian, uint64(*v))
	case uint:
		binary.Write(b, binary.BigEndian, uint64(v))
	case []int:
		binary.Write(b, binary.BigEndian, int64(len(v)))
		for _, i := range v {
			binary.Write(b, binary.BigEndian, uint64(i))
		}
	case []bool, []int8, []uint8, []int16, []uint16, []int32, []uint32, []int64, []uint64, []float32, []float64:
		binary.Write(b, binary.BigEndian, int64(reflect.ValueOf(a).Len()))
		binary.Write(b, binary.BigEndian, v)
	case *bool, bool,
		*int8, int8, *uint8, uint8,
		*int16, int16, *uint16, uint16,
		*int32, int32, *uint32, uint32,
		*int64, int64, *uint64, uint64,
		*float32, float32, *float64, float64:
		binary.Write(b, binary.BigEndian, v)
	default:
		// 尝试处理其他类型
		of := reflect.ValueOf(a)
		if !of.IsValid() {
			return false
		}
		switch of.Kind() {
		case reflect.Pointer, reflect.Invalid:
			return b.WriteType(of.Elem().Interface())
		}
		return false
	}
	return true
}

// ReadType 读取基础类型
func (b *Buffer) ReadType(a any) bool {
	switch v := a.(type) {
	case *[]bool, *[]int, *[]int8, *[]uint8, *[]int16, *[]uint16, *[]int32, *[]uint32, *[]int64, *[]uint64, *[]float32, *[]float64, *[]string:
		var le int64
		binary.Read(b, binary.BigEndian, &le)
		ov := reflect.ValueOf(a).Elem()
		of := reflect.MakeSlice(ov.Type(), int(le), int(le))
		for i := 0; i < int(le); i++ {
			b.ReadType(of.Index(i).Addr().Interface())
		}
		ov.Set(of)
	case *int:
		var a int64
		binary.Read(b, binary.BigEndian, &a)
		*v = int(a)
	case *bool, *int8, *uint8,
		*int16, *uint16, *int32, *uint32,
		*int64, *uint64, *float32, *float64:
		binary.Read(b, binary.BigEndian, v)
	case *string:
		var le int64
		binary.Read(b, binary.BigEndian, &le)
		*v = string(b.Next(int(le)))
	default:
		// 尝试处理其他类型
		of := reflect.ValueOf(a)
		if !of.IsValid() {
			return false
		}
		switch of.Kind() {
		case reflect.Interface:
			return b.ReadType(of.Elem())
		case reflect.Pointer:
			of := of.Elem()
			ov := reflect.New(of.Elem().Type())
			ok := b.ReadType(ov.Interface())
			of.Set(ov)
			return ok
		default:
			if of.CanAddr() {
				ov := reflect.New(of.Type())
				if b.ReadType(ov.Interface()) {
					of.Set(ov)
					return true
				}
			}
			return false
		}
	}
	return true
}

// WriteHex 写十六进制编码的字节数据
func (b *Buffer) WriteHex(src string) error {
	buufer, err := hex.DecodeString(src)
	if err == nil {
		b.Write(buufer)
	}
	return err
}

// ReadHex 将字节编码为十六进制
func (b *Buffer) ReadHex() string {
	return hex.EncodeToString(b.Bytes())
}

// Pool 缓冲池封装
type Pool struct {
	pool *sync.Pool
}

// NewPool 创建缓冲区
func NewPool(size int) Pool {
	if size == 0 {
		size = _size
	}
	return Pool{pool: &sync.Pool{
		New: func() any {
			return &Buffer{Buffer: bytes.NewBuffer(make([]byte, 0, _size))}
		},
	}}
}

// Get 得到对象
func (p Pool) Get() *Buffer {
	buf := p.pool.Get().(*Buffer)
	buf.Reset()
	buf.pool = p
	return buf
}

func (p Pool) put(buf *Buffer) {
	p.pool.Put(buf)
}
