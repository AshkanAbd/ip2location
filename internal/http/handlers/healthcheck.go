package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "healthy")
	}
}
