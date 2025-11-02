package handlers

import (
	"ip_location/internal/iptolocation/services"
	"net/http"

	"github.com/gin-gonic/gin"

	pkgMetrics "ip_location/pkg/metrics"
)

type HttpHandler struct {
	ipToLocation *services.IPToLocation
}

func NewHttpHandler(ipToLocation *services.IPToLocation) *HttpHandler {
	return &HttpHandler{
		ipToLocation: ipToLocation,
	}
}

func (h *HttpHandler) RegisterRoutes(engine *gin.Engine) {
	apiGp := engine.Group("/api")
	apiGp.GET("ip/:ip", h.GetIPInfo)
}

func (h *HttpHandler) GetIPInfo(ctx *gin.Context) {
	ip := ctx.Param("ip")
	if ip == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ip is required",
		})
		return
	}

	ipInfo, err := h.ipToLocation.GetIPInfo(ip)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	pkgMetrics.IncreaseCountryCount(ipInfo.Country)

	ctx.JSON(http.StatusOK, ipInfo)
}
