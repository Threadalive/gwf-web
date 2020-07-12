package middlewares

import (
	"gwf"
	"log"
	"time"
)

func Logger() gwf.HandleFunc {
	return func(c *gwf.Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
