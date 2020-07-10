package main

import (
	"gwf"
	"net/http"
)

func main() {
	r := gwf.New()

	//使用全局日志中间件
	r.Use(gwf.Logger())

	r.GET("/", func(c *gwf.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})

	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gwf.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

		v1.GET("/hello", func(c *gwf.Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2")
	{
		v3 := v2.Group("/v3")

		v3.GET("/", func(c *gwf.Context) {
			// expect /hello/dzx
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})

		v2.GET("/hello/:name", func(c *gwf.Context) {
			// expect /hello/dzx
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gwf.Context) {
			c.Json(http.StatusOK, gwf.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	r.Run(":8080")
}
