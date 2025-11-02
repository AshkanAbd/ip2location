package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	pkgMetrics "ip_location/pkg/metrics"
)

func Metrics() gin.HandlerFunc {
	h := promhttp.InstrumentMetricHandler(
		pkgMetrics.GetRegistry(), promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{}),
	)
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
