//测试
package main

import (
	"gwf"
	"net/http"
)

func main() {
	r := gwf.New()

	r.GET("/", func(c *gwf.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	r.GET("/hello", func(c *gwf.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.GET("/hello/:name", func(c *gwf.Context) {
		// expect /hello/geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *gwf.Context) {
		c.Json(http.StatusOK, gwf.H{"filepath": c.Param("filepath")})
	})

	r.Run(":9999")
}
