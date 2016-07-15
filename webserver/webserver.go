package webserver

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/colinyl/lib4go/logger"
	"github.com/colinyl/lib4go/utility"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Session string
	Address string
	Script  string
	Log     logger.ILogger
}

func NewContext(w http.ResponseWriter, r *http.Request, address string, script string) *Context {
	context := &Context{Writer: w, Request: r, Address: address, Script: script}
	context.Session = utility.GetSessionID()
	context.Log = logger.GetDeubgLogger(context.Session)
	return context

}

//WebHandler Web处理程序
type WebHandler struct {
	Path    string
	Script  string
	Method  string
	Handler func(*Context)
}

//WebServer WEB服务
type WebServer struct {
	routes     []WebHandler
	address    string
	loggerName string
	Log        logger.ILogger
}

//NewWebServer 创建WebServer服务
func NewWebServer(address string, loggerName string, handlers ...WebHandler) (server *WebServer) {
	server = &WebServer{routes: handlers, address: address, loggerName: loggerName}
	server.Log, _ = logger.Get(loggerName, true)
	return
}
func (h WebHandler) call(w http.ResponseWriter, r *http.Request) {
	context := NewContext(w, r, h.Path, h.Script)
	defer func() {
		if r := recover(); r != nil {
			context.Log.Fatal(r, string(debug.Stack()))
		}
	}()
	if strings.EqualFold(h.Method, "*") || strings.EqualFold(strings.ToLower(r.Method), strings.ToLower(h.Method)) {
		h.Handler(context)
		return
	}
	w.WriteHeader(404)
	w.Write([]byte("您访问的页面不存在"))
}

//Serve 启动WEB服务器
func (w *WebServer) Serve() {
	for _, handler := range w.routes {
		http.HandleFunc(handler.Path, handler.call)

	}
	err := http.ListenAndServe(w.address, nil)
	if err != nil {
		fmt.Println(err)
	}
}

//Stop 停止服务器
func (w *WebServer) Stop() {

}
