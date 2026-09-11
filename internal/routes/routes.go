package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kansaok/go-boilerplate/internal/config"
	"github.com/kansaok/go-boilerplate/internal/controller"
	"github.com/kansaok/go-boilerplate/internal/middleware"
	"github.com/kansaok/go-boilerplate/pkg/telemetry"
)

func SetupRoutes() *gin.Engine {
	r := gin.New()

	r.Use(gin.CustomRecovery(middleware.CustomRecovery))

	if err := r.SetTrustedProxies(nil); err != nil {
		panic(err)
	}

	r.Use(cors.New(config.CORSConfig()))

	r.Use(middleware.BodyLimitMiddleware(config.LoadConfig().SecurityConfig.MaxBodyBytes))

	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.TracingMiddleware())
	r.Use(middleware.RateLimitMiddleware(middleware.GlobalLimiter))

	r.Use(middleware.ValidateHost)
	r.Use(middleware.EnforceSSLRedirect)
	r.Use(middleware.SetSecurityHeaders)
	r.Use(middleware.SetCSRFHeaders)
	r.Use(middleware.SetSessionCookie)
	r.Use(middleware.SetXSSFilterHeader)
	r.Use(middleware.SetContentTypeNosniffHeader)

	r.GET("/metrics", middleware.MetricsProtection(), telemetry.PrometheusHandler())

	r.Use(middleware.CSRFValidation())

	api := r.Group("/api/v1")
	{
		authRoutes := api.Group("/auth")
		authRoutes.Use(middleware.RateLimitAuth())
		authRoutes.Use(middleware.BodyLimitMiddleware(1 << 20))
		AuthRoutes(authRoutes)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(config.LoadConfig().JWTConfig))
		{
			protected.GET("/me", controller.Me)
		}
	}

	return r
}