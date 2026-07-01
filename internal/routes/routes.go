package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kansaok/go-boilerplate/internal/config"
	"github.com/kansaok/go-boilerplate/internal/middleware"
	"github.com/kansaok/go-boilerplate/pkg/telemetry"
)

func SetupRoutes() *gin.Engine {
	// Inisialisasi Gin router
	r := gin.Default()

	// Setup CORS middleware
	r.Use(cors.New(config.CORSConfig()))

	// Add middleware
	r.Use(gin.CustomRecovery(middleware.CustomRecovery))
	r.Use(middleware.LoggingMiddleware())        // Custom logging middleware
	r.Use(middleware.TracingMiddleware())        // OpenTelemetry tracing middleware
	r.Use(middleware.RateLimitMiddleware(middleware.GlobalLimiter)) // Global rate limiting

	// Add security middleware
	r.Use(middleware.ValidateHost)
	r.Use(middleware.EnforceSSLRedirect)       		// Redirect HTTP to HTTPS
	r.Use(middleware.SetSecurityHeaders)        	// Security headers (CSP, HSTS, etc.)
	r.Use(middleware.SetCSRFHeaders)           		// Set CSRF cookie headers
	r.Use(middleware.SetSessionCookie)         		// Set secure session cookies
	r.Use(middleware.SetXSSFilterHeader)       		// XSS protection
	r.Use(middleware.SetContentTypeNosniffHeader) 	// Prevent content-type sniffing

	// Expose Prometheus metrics at /metrics endpoint
	r.GET("/metrics", telemetry.PrometheusHandler())

	api := r.Group("/api/v1")
	{
		// Auth routes with stricter rate limiting
		authRoutes := api.Group("/auth")
		authRoutes.Use(middleware.RateLimitAuth())
		AuthRoutes(authRoutes)
	}

	return r
}
