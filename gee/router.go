package gee

import (
	"net/http"
	"strings"
)

type router struct {
	// 该路径下所有的 trie 树的节点，包括 pattern, part, children
	roots map[string]*node
	// 方法加路径共同决定的对应的执行函数
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 处理 pattern，转换为 parts 获取路径信息
// eg: "/hello/world" => "", "hello", "world"
// 该函数支持处理 pattern = "*" 的情况
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		// 排除第一个为 "" 的情况
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 为路径添加对应的处理函数
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	// 处理 pattern 为parts
	parts := parsePattern(pattern)

	// key 为处理函数 handler 的键，key-handler 组成一对键值对
	key := method + "-" + pattern
	_, ok := r.roots[method]
	// 如果还没有创建 method 对应的 trie 树就创建一个
	if !ok {
		r.roots[method] = &node{}
	}
	// 插入对应节点并记录处理函数
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// 获取对应的节点
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	// 处理 path 获得要寻找的 searchPath
	searchParts := parsePattern(path)

	// 初始定义 params
	params := make(map[string]string)
	// 获取 method 对应的 trie 树
	root, ok := r.roots[method]
	// 没有该方法的 trie 树
	if !ok {
		return nil, nil
	}

	// 查询目标节点
	n := root.search(searchParts, 0)

	// 节点存在
	if n != nil {
		// 因为涉及 ":" 和 "*" 所以不能用参数中的 Path
		// 将返回的节点路径 pattern 转换为 parts
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			// 涉及模糊判断，将对应的模糊判断的实际参数记录
			// eg：part = ":leng", searchParts[index] = "c" => params["leng"] = "c"
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			// 涉及通配判断，将通配判断的实际参数记录
			// eg: part = "*filepath", searchParts[index] = "src", searchParts[index+1] = "index.html"
			// => params["filepath"] = "src/index.html"
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	// 节点不存在
	return nil, nil
}

// 获取某一方法下所有的实际节点
func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

// 封装处理函数
func (r *router) handle(c *Context) {
	// 获取实际对应的节点
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		// 记录涉及到的 ":" "*" 参数
		c.Params = params
		// 获取实际调用的 handler 的键 key
		key := c.Method + "-" + n.pattern
		// 运行 handler
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
