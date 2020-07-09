package gwf

import (
	"fmt"
	"net/http"
)

//定义请求的处理函数主体
type HandleFunc func(w http.ResponseWriter, r *http.Request)

//实现ServeHTTP的实例
type Engine struct {
	//路由和处理函数的映射
	router map[string]HandleFunc
}

//创建gwf实例的函数,返回一个实例引用
func New() *Engine {
	return &Engine{map[string]HandleFunc{}}
}

func (engine *Engine) addRoute(method string, pattern string, handleFunc HandleFunc) {
	key := method + "-" + pattern
	engine.router[key] = handleFunc
}

//定义GET方法接收
func (engine *Engine) GET(pattern string, handleFunc HandleFunc) {
	engine.addRoute("GET", pattern, handleFunc)
}

//定义POST方法接收
func (engine *Engine) POST(pattern string, handleFunc HandleFunc) {
	engine.addRoute("POST", pattern, handleFunc)
}

//定义GET方法接收
func (engine *Engine) Run(addr string) {
	http.ListenAndServe(addr, engine)
}

//实现接口
func (engine *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key := request.Method + "-" + request.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(writer, request)
	} else {
		fmt.Fprintf(writer, "404 NOT FOUND: %s\n", request.URL)
	}
}
