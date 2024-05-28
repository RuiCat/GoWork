package logs

import (
	"fmt"
	"os"
)

// LogLevel 日志等级
type LogLevel uint8

// 默认日志等级
var (
	DebugLevel   = NewLevel(0x0, "DebugLevel", nil)   // "Debug.log")
	InfoLevel    = NewLevel(0x1, "InfoLevel", nil)    //"Info.log")
	WarnLevel    = NewLevel(0x2, "WarnLevel", nil)    //"Warn.log")
	ErrorLevel   = NewLevel(0x3, "ErrorLevel", nil)   // "Error.log")
	FatalLevel   = NewLevel(0x4, "FatalLevel", nil)   // "Fatal.log")
	UnknownLevel = NewLevel(0x5, "UnknownLevel", nil) // "Unknown.log")
)

// LogInfo 日志信息
type LogInfo struct {
	Name    string
	Level   LogLevel
	LogFile *os.File
}

// NewLevel 创建新日志等级
func NewLevel(Level LogLevel, Name string, File *os.File) LogLevel {
	if _, ok := Std.LogMap.Load(Level); !ok {
		if File == nil {
			File = Std.LogOut
		}
		Std.LogMap.Store(Level, &LogInfo{Name: Name, Level: Level, LogFile: File})
		return Level
	}
	panic(fmt.Sprintf("函数 NewLevel(Level:%d Name:%s File%s) 重复创建", Level, Name, File.Name()))
}

// NewCode 创建新日志代码
func (Level LogLevel) NewCode(Code LogCode) LogCode {
	if _, ok := Std.LogCode.Load(Code); !ok {
		Std.LogCode.Store(Code, Level)
		return Code
	}
	panic(fmt.Sprintf("函数 NewCode(Code:%d,Level:%d) 重复创建", Code, Level))
}

// Info 日志信息
func (level LogLevel) Info() *LogInfo {
	if info, ok := Std.LogMap.Load(level); ok {
		return info.(*LogInfo)
	}
	return nil
}

// String 等级标志
func (level LogLevel) String() string {
	if info, ok := Std.LogMap.Load(level); ok {
		return info.(*LogInfo).Name
	}
	return fmt.Sprintf("LogLevel(%d?)", level)
}

// LogCode 日志代码
type LogCode uint64

// Level 日志等级
func (code LogCode) Level() LogLevel {
	if level, ok := Std.LogCode.Load(code); ok {
		return level.(LogLevel)
	}
	return UnknownLevel
}

// Info 日志信息
func (code LogCode) Info() *LogInfo {
	return code.Level().Info()
}

// String 代码标志
func (code LogCode) String() string {
	return fmt.Sprintf("Code:%d Level:%s", code, code.Level().String())
}

// Print 写入信息
func (code LogCode) Print(a ...any) {
	code.Info().LogFile.WriteString(fmt.Sprintf("%s,Details:%v\n", code.String(), a))
}

// Error 输出错误
func (code LogCode) Error(err error) {
	code.Print(err.Error())
}

// IsError 检测错误
func (code LogCode) IsError(err error, info ...string) {
	if err != nil {
		code.Print(fmt.Sprintf("%s,info:%v\n", err.Error(), info))
	}
}

// Defer 拦截错误
func (code LogCode) Defer() {
	if err := recover(); err != nil {
		code.Print(err)
	}
}

// New 构建错误
func (code LogCode) New(a any, Info ...any) error {
	return &LogDetails{Code: code, Details: a, Info: Info}
}

// LogDetails 日志记录
type LogDetails struct {
	Code    LogCode // 日志代码
	Details any     // 日志内容
	Info    []any   // 附加信息
}

func (log LogDetails) Error() string {
	buffer := Std.Pool.Get()
	defer buffer.Free()
	buffer.WriteString("Code:[")
	buffer.WriteTypeString(log.Code)
run:
	if len(log.Info) > 0 {
		buffer.WriteString(",Info:[")
		for _, v := range log.Info {
			buffer.WriteTypeString(v)
			buffer.WriteString(",")
		}
		buffer.Truncate(buffer.Len() - 1)
		buffer.WriteString("]")
	}
	switch v := log.Details.(type) {
	case *LogDetails:
		log = *v
	case LogDetails:
		log = v
	default:
		goto exit
	}
	buffer.WriteString("->")
	buffer.WriteTypeString(log.Code)
	goto run
exit:
	buffer.WriteString("],Details:")
	buffer.WriteTypeString(log.Details)
	return buffer.String()
}

// IfPrint 输出错误/日志信息
func IfPrint(is bool, err error) {
	if is && err != nil {
		switch v := err.(type) {
		case LogDetails:
			v.Code.Error(v)
		case *LogDetails:
			v.Code.Error(v)
		default:
			print(v)
		}
	}
}

// Print 输出错误/日志信息
func Print(err error) {
	if err != nil {
		switch v := err.(type) {
		case LogDetails:
			v.Code.Error(v)
		case *LogDetails:
			v.Code.Error(v)
		default:
			print(v)
		}
	}
}

// Panic 输出错误/日志信息
func Panic(err error) {
	if err != nil {
		switch v := err.(type) {
		case LogDetails:
			v.Code.Error(v)
		case *LogDetails:
			v.Code.Error(v)
		default:
			panic(v)
		}
	}
}
