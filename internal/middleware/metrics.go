package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func MetricsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        promhttp.Handler().ServeHTTP(c.Writer, c.Request)
    }
}
