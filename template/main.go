package main

import (
	"example/middlewares"
	"fmt"
	"gwf"
	"html/template"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age  int8
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gwf.New()
	r.Use(middlewares.Logger())

	r.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	//加载模板目录块
	r.LoadHTMLGlob("templates/*")
	//将用户url中访问资源的路径/assets与本地/static目录绑定
	r.Static("/assets", "./static")

	stu1 := &student{Name: "dzx", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}

	r.GET("/", func(c *gwf.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})

	r.GET("/student", func(c *gwf.Context) {
		c.HTML(http.StatusOK, "stu_arr.tmpl", gwf.H{
			"title":  "gwf",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	r.GET("/date", func(c *gwf.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gwf.H{
			"title": "gwf",
			"now":   time.Date(2020, 7, 12, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":8080")
}
