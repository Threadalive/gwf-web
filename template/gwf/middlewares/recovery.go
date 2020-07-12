package middlewares

import (
	"fmt"
	"gwf"
	"log"
	"net/http"
	"runtime"
	"strings"
)

//打印堆栈日志函数
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

//用来宕机恢复的中间件
func Recovery() gwf.HandleFunc {
	return func(c *gwf.Context) {
		defer func() {
			//宕机恢复
			if err := recover(); err != nil {
				msg := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(msg))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		c.Next()
	}
}
