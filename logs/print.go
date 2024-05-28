package logs

import (
	"fmt"
	"io"
	"os"
	"reflect"
)

// PrintAny 输出
func PrintAny(x any, tier int) {
	p := printer{
		output: os.Stdout,
		ptrmap: make(map[any]int),
		tier:   tier,
		last:   '\n', // force printing of line number on first line
	}
	p.print(reflect.ValueOf(x))
}

// PrintAnyw 输出
func PrintAnyw(w io.Writer, x any, tier int) {
	p := printer{
		output: w,
		ptrmap: make(map[any]int),
		tier:   tier,
		last:   '\n', // force printing of line number on first line
	}
	p.print(reflect.ValueOf(x))
}

// printer 打印输出
type printer struct {
	output io.Writer
	ptrmap map[any]int // *T -> line number
	indent int         // current indentation level
	last   byte        // the last byte processed by Write
	line   int         // current line number
	tier   int
}

var indent = []byte(".  ")

func (p *printer) Write(data []byte) (n int, err error) {
	var m int
	for i, b := range data {
		// invariant: data[0:n] has been written
		if b == '\n' {
			m, err = p.output.Write(data[n : i+1])
			n += m
			if err != nil {
				return
			}
			p.line++
		} else if p.last == '\n' {
			_, err = fmt.Fprintf(p.output, "%6d  ", p.line)
			if err != nil {
				return
			}
			for j := p.indent; j > 0; j-- {
				_, err = p.output.Write(indent)
				if err != nil {
					return
				}
			}
		}
		p.last = b
	}
	if len(data) > n {
		m, err = p.output.Write(data[n:])
		n += m
	}
	return
}

type localError struct{ err error }

func (p *printer) Printf(format string, args ...any) {
	p.printf(format, args...)
}
func (p *printer) printf(format string, args ...any) {
	if _, err := fmt.Fprintf(p, format, args...); err != nil {
		panic(localError{err})
	}
}

func (p *printer) print(x reflect.Value) {
	if x.IsZero() {
		p.printf("nil")
		return
	}
	if x.CanAddr() {
		ptr := x.UnsafeAddr()
		if line, exists := p.ptrmap[ptr]; exists {
			p.printf("(obj @ %d)", line)
			return
		} else {
			p.ptrmap[ptr] = p.line
		}
	}
	if p.tier != -1 && p.indent > p.tier {
		p.printf("(Add @ 0x%x) %#+v", x.UnsafeAddr(), x)
		return
	}
	switch x.Kind() {
	case reflect.Interface:
		if x.CanInterface() {
			p.print(x.Elem())
		} else {
			p.printf("%v", x.Elem())
		}
	case reflect.Map:
		p.printf("%s (len = %d) {", x.Type(), x.Len())
		if x.Len() > 0 {
			p.indent++
			p.printf("\n")
			for _, key := range x.MapKeys() {
				p.print(key)
				p.printf(": ")
				p.print(x.MapIndex(key))
				p.printf("\n")
			}
			p.indent--
		}
		p.printf("}")
	case reflect.Pointer:
		p.printf("*")
		var ptr any
		if x.CanInterface() {
			ptr = x.Interface()
		} else {
			if x.CanAddr() {
				ptr = x.Addr()
			} else {
				ptr = x
			}
		}
		if line, exists := p.ptrmap[ptr]; exists {
			p.printf("(obj @ %d)", line)
		} else {
			p.ptrmap[ptr] = p.line
			p.print(x.Elem())
		}
	case reflect.Array:
		p.printf("%s {", x.Type())
		if x.Len() > 0 {
			p.indent++
			p.printf("\n")
			for i, n := 0, x.Len(); i < n; i++ {
				p.printf("%d: ", i)
				p.print(x.Index(i))
				p.printf("\n")
			}
			p.indent--
		}
		p.printf("}")
	case reflect.Slice:
		if x.CanInterface() {
			if s, ok := x.Interface().([]byte); ok {
				p.printf("%#q", s)
				return
			}
		}
		p.printf("%s (len = %d) {", x.Type(), x.Len())
		if x.Len() > 0 {
			p.indent++
			p.printf("\n")
			for i, n := 0, x.Len(); i < n; i++ {
				p.printf("%d: ", i)
				p.print(x.Index(i))
				p.printf("\n")
			}
			p.indent--
		}
		p.printf("}")
	case reflect.Struct:
		t := x.Type()
		p.printf("%s {", t)
		p.indent++
		first := true
		for i, n := 0, t.NumField(); i < n; i++ {
			name := t.Field(i).Name
			value := x.Field(i)
			if first {
				p.printf("\n")
				first = false
			}
			p.printf("%s: ", name)
			p.print(value)
			p.printf("\n")
		}
		p.indent--
		p.printf("}")
	default:

		if x.CanInterface() {
			fmt.Println("遇到错误:", x, "<-这是啥玩意")
			v := x.Interface()
			switch v := v.(type) {
			case string:
				p.printf("%s", v)
				return
			}
		}
		p.printf("%+v", x)
	}
}
