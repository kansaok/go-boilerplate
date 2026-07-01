package telemetry

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var excludedRoutes = map[string]bool{
    "/product": true,
    "/marl":    true,
    "/inces":   true,
}

var httpRequests = prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Total number of HTTP requests",
    },
    []string{"method", "path"},
)

func InitPrometheus() {
    prometheus.MustRegister(httpRequests)
}

func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
		// Check if the current route is excluded
        if _, ok := excludedRoutes[c.FullPath()]; ok {
            c.Next() // Skip middleware for excluded routes
            return
        }
        path := c.FullPath()
        method := c.Request.Method

        // Count the request
        httpRequests.WithLabelValues(method, path).Inc()

        c.Next()
    }
}

func PrometheusHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        promhttp.Handler().ServeHTTP(c.Writer, c.Request)
    }
}
