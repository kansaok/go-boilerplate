package config

import (
	"os"
	"strings"

	"github.com/gin-contrib/cors"
)

func CORSConfig() cors.Config {
	origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
	if len(origins) == 0 || (len(origins) == 1 && origins[0] == "") {
		origins = []string{"http://localhost:3000"}
	}

	return cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * 60 * 60,
	}
}
