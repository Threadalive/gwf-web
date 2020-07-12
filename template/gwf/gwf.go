package gwf

import (
	"gwf/middlewares"
	"html/template"
	"log"
	"net/http"
	"path"
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
		//html模板
		htmlTemplates *template.Template
		//用户可自定义的渲染函数集
		funcMap template.FuncMap
	}
)

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

//加载html目录文件
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandleFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// serve static files
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	group.GET(urlPattern, handler)
}

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

//默认创建实例中使用日志记录和宕机恢复中间件
func Default() *Engine {
	engine := New()
	engine.Use(middlewares.Logger(), middlewares.Recovery())

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
	//设置context中的engine为当前engine
	context.engine = engine

	engine.router.handle(context)
}
