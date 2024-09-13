package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/kansaok/go-boilerplate/docs"
	"github.com/kansaok/go-boilerplate/internal/config"
	"github.com/kansaok/go-boilerplate/internal/middleware"
	"github.com/kansaok/go-boilerplate/pkg/telemetry"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
    // r.Use(middleware.PrometheusMiddleware())     // Prometheus metrics middleware

	// Add security middleware
	r.Use(middleware.ValidateHost)
	r.Use(middleware.EnforceSSLRedirect)       		// Redirect HTTP to HTTPS
	r.Use(middleware.SetCSRFHeaders)           		// Set CSRF cookie headers
	r.Use(middleware.SetSessionCookie)         		// Set secure session cookies
	r.Use(middleware.SetXSSFilterHeader)       		// XSS protection
	r.Use(middleware.SetContentTypeNosniffHeader) 	// Prevent content-type sniffing

	// Expose Prometheus metrics at /metrics endpoint
	r.GET("/metrics", telemetry.PrometheusHandler())
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/api/v1")
    {
        AuthRoutes(api)
    }

	return r
}
