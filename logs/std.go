package logs

import (
	"os"
	"sync"
)

// Std 全局函数
var Std = struct {
	LogOut  *os.File
	LogMap  sync.Map
	LogCode sync.Map
	Exit    chan struct{}
	Pool    Pool
}{Exit: make(chan struct{}), LogOut: os.Stdout, Pool: NewPool(1024)}

// Exit 结束程序
func Exit() {
	close(Std.Exit)
}

// IsExit 程序是否结束
func IsExit() <-chan struct{} {
	return Std.Exit
}

func init() {
	// 更改默认输出
	os.Stdout = InfoLevel.Info().LogFile
	os.Stderr = ErrorLevel.Info().LogFile
}

// Close 用于处理代码中存在多次关联的 Close 过程
func Close(close *func()) (func(), *bool, func(fn any)) {
	var i int
	var ok bool
	add := make([]any, 0, 10)
	(*close) = func() {
		close = nil
		for ; i < len(add); i++ {
			switch fn := add[i].(type) {
			case func():
				fn()
			case func() error:
				fn()
			}
		}
	}
	return func() {
			if !ok {
				if close != nil {
					(*close)()
				}
			}
		}, &ok, func(fn any) {
			add = append(add, fn)
		}
}
