package gwf

import (
	"log"
	"net/http"
	"strings"
)

//定义请求的处理函数主体
type HandleFunc func(*Context)

type (
	//路由分组
	RouterGroup struct {
		prefix string
		//中间件支持
		middlewares []HandleFunc
		//parent *RouterGroup
		engine *Engine
	}
	//实现ServeHTTP的实例
	Engine struct {
		*RouterGroup
		//路由和处理函数的映射
		router *router
		//存储所有路由分组
		groups []*RouterGroup
	}
)

//使用中间件
func (group *RouterGroup) Use(middlewares ...HandleFunc) {
	group.middlewares = append(group.middlewares, middlewares...)

}

//创建gwf实例的函数,返回一个实例引用
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}

	return engine
}

//新建路由分组
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		//parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

//注册路由
func (group *RouterGroup) addRoute(method string, comp string, handleFunc HandleFunc) {
	//完整路由
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)

	group.engine.router.addRoute(method, pattern, handleFunc)
}

//定义GET方法接收
func (group *RouterGroup) GET(pattern string, handleFunc HandleFunc) {
	group.addRoute("GET", pattern, handleFunc)
}

//定义POST方法接收
func (group *RouterGroup) POST(pattern string, handleFunc HandleFunc) {
	group.addRoute("POST", pattern, handleFunc)
}

//定义GET方法接收
func (engine *Engine) Run(addr string) {
	http.ListenAndServe(addr, engine)
}

//实现接口
func (engine *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var middlewares []HandleFunc

	for _, group := range engine.groups {
		//若请求路由在当前分组中,则添加当前分组的中间件
		if strings.HasPrefix(request.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	context := NewContext(writer, request)
	context.handler = middlewares

	engine.router.handle(context)
}
