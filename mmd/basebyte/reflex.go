package basebyte

import (
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

var IntSize int = 32 << (^uint(0) >> 63) / 8

func intDataSize(kind reflect.Kind) int {
	switch kind {
	case reflect.Bool, reflect.Int8, reflect.Uint8:
		return 1
	case reflect.Int, reflect.Ptr, reflect.UnsafePointer:
		return IntSize
	case reflect.Int16, reflect.Uint16:
		return 2
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		return 4
	case reflect.Int64, reflect.Uint64, reflect.Float64, reflect.Complex64, reflect.Map:
		return 8
	case reflect.Complex128:
		return 16
	}
	panic(fmt.Sprintf("复合类型数据无法读取: %v", kind))
}
func getGyte(size int, pointer uintptr, w io.Writer) {
	b := make([]byte, size)
	for i := 0; i < size; i++ {
		b[i] = *(*byte)(unsafe.Pointer(pointer + uintptr(i)))
	}
	w.Write(b)
}

// KindField 子类型
type KindField struct {
	Offset uintptr // 偏移
	Value  *Kind   // 子类型
}

// Kind 类型
type Kind struct {
	Type  reflect.Type // 原始类型
	Value interface{}  // 附加数据
	Node  []KindField  // 由组合数据构成的链表
}

// Get 得到数据
func (kind *Kind) Get(pointer uintptr, w io.Writer) {
	if pointer == 0 {
		return
	}
	switch kind.Type.Kind() {
	case reflect.Struct:
		le := len(kind.Node)
		for i := 0; i < le; i++ {
			kind.Node[i].Value.Get(pointer+kind.Node[i].Offset, w)
		}
	case reflect.String:
		le := *(*int)(unsafe.Pointer(pointer + uintptr(IntSize)))
		getGyte(IntSize, pointer+uintptr(IntSize), w)
		if le > 0 {
			ptr := *(*uintptr)(unsafe.Pointer(pointer))
			getGyte(le, ptr, w)
		}
	case reflect.Slice:
		header := *(*reflect.SliceHeader)((unsafe.Pointer)(pointer))
		pointer = header.Data
		kind.Value = header.Len
		fallthrough
	case reflect.Array:
		// 写入数组数量
		le := kind.Value.(int)
		getGyte(IntSize, (uintptr)(unsafe.Pointer(&le)), w)
		// 读取数据元素
		offset := kind.Node[0].Offset
		ptrKind := kind.Node[0].Value
		// 计算数据偏移
		for i := uintptr(0); i < uintptr(le); i++ {
			ptrKind.Get(pointer+(i*offset), w)
		}
	case reflect.Map:
		// 构建类型
		ptrKey := kind.Node[0].Value
		ptrValue := kind.Node[1].Value
		newKey := reflect.New(ptrKey.Type)
		newValue := reflect.New(ptrValue.Type)
		// 得到数据指针
		keyPtr := newKey.Pointer()
		valuePtr := newValue.Pointer()
		newKey = newKey.Elem()
		newValue = newValue.Elem()
		// 指针还原构建
		ptrMap := reflect.NewAt(kind.Type, unsafe.Pointer(pointer)).Elem()
		// 数据数量
		le := ptrMap.Len()
		getGyte(IntSize, (uintptr)(unsafe.Pointer(&le)), w)
		// ptrMap
		iter := ptrMap.MapRange()
		for iter.Next() {
			newKey.Set(iter.Key())
			newValue.Set(iter.Value())
			ptrKey.Get(keyPtr, w)
			ptrValue.Get(valuePtr, w)
		}
	case reflect.Ptr, reflect.UnsafePointer:
		kind.Node[0].Value.Get(*(*uintptr)(unsafe.Pointer(pointer)), w)
	default:
		getGyte(kind.Value.(int), pointer, w)
	}
}

// Analysis 分析结构体
func Analysis(in interface{}) *Kind {
	return analysis(reflect.TypeOf(in))
}

func analysis(value reflect.Type) (kind *Kind) {
	// 构建与原始数据无关的类型
	kind = &Kind{
		Type: value,
	}
	// 对数据进行处理
	switch kind.Type.Kind() {
	case reflect.Struct: // 组合类型的结构需要对结构体的子结构进行递归处理
		for i, num := 0, value.NumField(); i < num; i++ {
			field := value.Field(i)
			kind.Node = append(kind.Node, KindField{
				Offset: field.Offset,
				Value:  analysis(field.Type),
			})
		}
	case reflect.Array:
		kind.Value = value.Len()
		fallthrough
	case reflect.Slice, reflect.Ptr, reflect.UnsafePointer:
		kind.Node = append(kind.Node, KindField{
			Value:  analysis(value.Elem()),
			Offset: value.Elem().Size(),
		})
	case reflect.Map:
		kind.Node = append(kind.Node, KindField{
			Value: analysis(kind.Type.Key()),
		})
		kind.Node = append(kind.Node, KindField{
			Value: analysis(kind.Type.Elem()),
		})
	case reflect.Chan, reflect.Func, reflect.Interface:
		panic(fmt.Sprintf("无法序列化类型: %v", kind))
	case reflect.String:
	default:
		kind.Value = intDataSize(kind.Type.Kind())
	}
	return kind
}
