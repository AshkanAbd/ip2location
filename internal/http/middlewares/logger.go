package middlewares

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	pkgLog "ip_location/pkg/logger"
	pkgMetrics "ip_location/pkg/metrics"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		path := c.FullPath()
		if strings.HasPrefix(path, "/api/") {
			pkgMetrics.IncreaseHttpStatus(fmt.Sprintf("%d", status))
		}

		switch {
		case status < 400:
			pkgLog.Debug("%d %s %s %.3f ms", status, c.Request.Method, path, float64(latency.Microseconds())/1000)
		case status >= 400 && status < 500:
			pkgLog.Info("%d %s %s %.3f ms", status, c.Request.Method, path, float64(latency.Microseconds())/1000)
		case status >= 500:
			pkgLog.Warn("%d %s %s %.3f ms", status, c.Request.Method, path, float64(latency.Microseconds())/1000)
		}
	}
}
