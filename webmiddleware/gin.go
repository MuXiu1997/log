package webmiddleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muxiu1997/log"
)

func GinLogger() gin.HandlerFunc {
	gin.Logger()
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		timeStamp := time.Now()
		latency := timeStamp.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		msg := fmt.Sprintf("%3d | %13v | %15s | %-7s  %#v",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)

		entry := log.Logger().WithFields(map[string]interface{}{
			"latency":    latency,
			"clientIP":   clientIP,
			"method":     method,
			"statusCode": statusCode,
			"bodySize":   bodySize,
		})

		if 500 <= statusCode {
			entry.Errorln(msg)
		} else {
			entry.Infoln(msg)
		}
	}
}
