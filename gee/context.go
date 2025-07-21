package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin object

	// Writer 是响应的重要内容，包括响应头，状态码，和主体
	Writer http.ResponseWriter
	// Req 是从客户端获取到的请求文件，包括url路由和请求方式
	Req *http.Request

	// request info

	// 从 Req 中获取的 url path
	Path string
	// 从 Req 中获取的 method
	Method string
	// 涉及到的 ":" "*" 的对应参数
	Params map[string]string

	// response info

	// 状态码
	StatusCode int

	// middleware

	// 记录所有处理函数，
	// 其内部的所有待处理函数会依次执行
	handlers []HandlerFunc
	// 记录已执行处理函数，
	// 由于执行处理函数时可能会出现再次调用 c.Next() 方法，
	// 所以需要将 index 设置在 Context 内
	index int
}

// 根据 http.Reaponse 和 *http.Request 创建一个 Context
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// 执行所有处理函数
func (c *Context) Next() {
	// 获取下一个待执行函数
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		// 执行
		c.handlers[c.index](c)
	}
}

// waiting...
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// 根据需要获得的参数名来获取对应的参数
// ":" => 路径中的实际路径名
// "*" => 通配符后的所有内容
func (c *Context) Param(key string) string {
	return c.Params[key]
}

// 获取在 url 和 POST 请求体中变量第一次出现时对应的参数
// 通常在 url 和 POST 请求体中同时存在同名变量时，会获取到 url 中的参数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 仅获取 url 中的参数
// eg: url = "http://localhost:8089/?uid=123"
// Quary return "123"
func (c *Context) Quary(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 设置状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置相应头内容
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 设置响应体内容-String
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Context-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 设置响应体内容-JSON
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Context-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 设置响应体内容-Data
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// 设置响应体内容-HTML
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Context-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
