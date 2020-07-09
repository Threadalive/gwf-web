package gwf

import (
	"log"
	"net/http"
	"strings"
)

//路由表
type router struct {
	roots    map[string]*node
	handlers map[string]HandleFunc
}

func newRouter() *router {
	return &router{
		roots:    map[string]*node{},
		handlers: make(map[string]HandleFunc),
	}
}

//根据"/"分割路由，返回路由各个部分的数组
func parsePattern(pattern string) []string {
	temps := strings.Split(pattern, "/")

	parts := make([]string, 0)

	for _, part := range temps {
		if part != "" {
			parts = append(parts, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return parts
}

//注册路由
func (r *router) addRoute(method string, pattern string, handler HandleFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	//解析获取路由各部分
	parts := parsePattern(pattern)

	key := method + "-" + pattern

	_, ok := r.roots[method]
	//若该方法暂时不存在路由根节点，则新建一个节点
	if !ok {
		r.roots[method] = &node{}
	}
	//注册该路由，插入节点
	r.roots[method].insert(pattern, parts, 0)

	r.handlers[key] = handler
}

//获取用户查询的路由
func (r *router) getRouter(method string, path string) (*node, map[string]string) {
	//分割用户查询的路径
	searchParts := parsePattern(path)

	params := make(map[string]string)

	root, ok := r.roots[method]
	//若该方法对应的路由未注册，返回空
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				//拼接后续路由
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

//处理函数
func (r *router) handle(c *Context) {
	n, params := r.getRouter(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		//执行对应函数
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
