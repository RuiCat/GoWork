package webui

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"webui/websocket"
)

//go:embed PetUI.js
var webPetui []byte
var webPetuiLength string = strconv.FormatInt(int64(len(webPetui)), 10)
var upgrader = websocket.Upgrader{}

// EventFn 函数接口
type EventFn func(arg any)

// Event 事件接口
type Event interface {
	Close()
	CallFunc(name string, value any)
	SetValue(name string, value any)
	GetValue(name string) (value any)
	DeleteValue(name string)
	SetFunc(name string, fn EventFn)
	DeleteFunc(name string)
}

// value 事件值
type value struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value any    `json:"value"`
}
type event struct {
	conn      *websocket.Conn
	isClose   bool
	funclist  sync.Map
	valueList sync.Map
}

func (eve *event) Close() {
	eve.isClose = false
}
func (eve *event) CallFunc(name string, val any) {
	eve.conn.WriteJSON(&value{
		Type:  "CallFunc",
		Name:  name,
		Value: val,
	})
}
func (eve *event) SetValue(name string, val any) {
	eve.valueList.Store(name, val)
	eve.conn.WriteJSON(&value{
		Type:  "SetValue",
		Name:  name,
		Value: val,
	})
}
func (eve *event) GetValue(name string) (val any) {
	if v, ok := eve.valueList.Load(name); ok {
		return v
	}
	return nil
}
func (eve *event) DeleteValue(name string) {
	eve.valueList.Delete(name)
	eve.conn.WriteJSON(&value{
		Type: "DeleteValue",
		Name: name,
	})
}
func (eve *event) SetFunc(name string, fn EventFn) {
	eve.funclist.Store(name, fn)
	eve.conn.WriteJSON(&value{
		Type: "SetFunc",
		Name: name,
	})
}
func (eve *event) DeleteFunc(name string) {
	eve.funclist.Delete(name)
	eve.conn.WriteJSON(&value{
		Type: "DeleteFunc",
		Name: name,
	})
}

// PetUI 创建UI框架
func PetUI(pattern string, fn func(r *http.Request, event Event)) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		// 升级协议
		conn, _ := upgrader.Upgrade(w, r, nil)
		defer conn.Close()
		// 事件结构体
		eve := &event{conn: conn, isClose: true}
		fn(r, eve)
		val := &value{}
		val.Type = "onopen"
		eve.conn.WriteJSON(val)
		for eve.isClose {
			messageType, message, _ := conn.ReadMessage()
			if messageType == -1 {
				eve.isClose = false
				return
			}
			json.Unmarshal(message, val)
			// 处理事件
			switch val.Type {
			case "SetValue":
				eve.valueList.Store(val.Name, val.Value)
			case "DeleteFunc":
				eve.funclist.Delete(val.Name)
			case "DeleteValue":
				eve.valueList.Delete(val.Name)
			case "CallFunc":
				if fn, ok := eve.funclist.Load(val.Name); ok {
					(fn.(EventFn))(val.Value)
				}
			default:
			}
		}
	})
	// 打包JS库
	mux.HandleFunc(pattern+".js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript;charset=utf-8")
		w.Header().Set("Content-Length", webPetuiLength)
		w.Write(webPetui)
	})
	return mux
}
