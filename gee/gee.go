package gee

import (
	"net/http"
)

// 定义 HandlerFunc 为以 *Context 为参数的函数
type HandlerFunc func(*Context)

type Engine struct {
	route *router
}

func New() *Engine {
	return &Engine{route: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.route.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) error {
	e := http.ListenAndServe(addr, engine)
	return e
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.route.handle(c)
}
