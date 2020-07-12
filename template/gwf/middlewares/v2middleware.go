package middlewares

import (
	"gwf"
	"log"
	"time"
)

//进在v2分组下执行的中间件
func OnlyForV2() gwf.HandleFunc {
	return func(c *gwf.Context) {
		t := time.Now()
		//若服务器发生错误
		c.Fail(500, "Internal Server Error")

		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
