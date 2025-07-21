package gee

import (
	"net/http"
)

// 定义 HandlerFunc 为以 *Context 为参数的函数
type HandlerFunc func(*Context)

// 路由组实现分组功能
type RouterGroup struct {
	// 组名，
	// 如我有一些在 user 下的路由，
	// userGroup := server.Group("/user")
	// 则此时一个 RouterGroup 的实例 userGroup 的 prefix = "/user"
	// 或者在 "/" 的路由，
	// root := server.Group("/")
	// root.prefix = "/"
	prefix string
	// 中间件，在同一个路由组的下，均会触发该组的中间件函数
	middlewares []HandlerFunc
	// 父节点路由组
	parent *RouterGroup
	// 全局的 engine 指针，确保所有资源由 engine 控制
	engine *Engine
}

type Engine struct {
	// 接口继承自 RouterGroup
	*RouterGroup
	router *router
	groups []*RouterGroup
}

// 在主函数中定义的 server := gee.New()
// 其中 server 是 Engine 同时也是继承自 RouterGroup
// Engine 唯一存在，负责管理所有资源
// RouterGroup 来辨析所有路由路径
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		// 链接之前路由组的组名
		// 实现路由组的嵌套
		// 如：
		// user := server.Group("/user")
		// login := user.Group("/login")
		// 则 RouterGroup 实例 login.prefix = "/user/login"
		prefix: group.prefix + prefix,
		// 记录该路由组节点的父节点
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	// 分组后，在添加新路由时，只需提供后续路由位置
	// 但在 trie 树中需要存储的是实际的路由字符串
	// 故需要通过 group.prefix 来补全
	pattern := group.prefix + comp
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) error {
	e := http.ListenAndServe(addr, engine)
	return e
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
