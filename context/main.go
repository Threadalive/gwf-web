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
		// expect /hello?name=dzx
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *gwf.Context) {
		c.Json(http.StatusOK, gwf.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":8080")
}
