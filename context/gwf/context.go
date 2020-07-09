//提供了访问Query和PostForm参数的方法。
//提供了快速构造String/Data/JSON/HTML响应的方法。
package gwf

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	// response info
	StatusCode int
}

//新建请求的context实例
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
	}
}

//POST方法解析参数
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

//GET方法解析获取参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

//设置请求状态
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

//设置返回头信息
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

//json格式返回
func (c *Context) Json(code int, values interface{}) {
	c.SetHeader("Content-Type", "text/json")
	//设置状态码
	c.Status(code)
	//新建一个json编码器
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(values); err != nil {
		http.Error(c.Writer, err.Error(), 500)
		//panic(err)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
