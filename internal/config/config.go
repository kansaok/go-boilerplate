package config

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

// AppConfig stores all application configurations.
type AppConfig struct {
	DatabaseConfig 	*Config
	JWTConfig      	*JWTConfig
	CORSConfig     	cors.Config
	SecurityConfig	SecurityConfig
}

// LoadConfig loads all application configurations.
func LoadConfig() *AppConfig {
	// Create file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Return all configurations in a single struct.
	return &AppConfig{
		DatabaseConfig: LoadDatabaseConfig(),
		JWTConfig:      LoadJWTConfig(),
		CORSConfig:     CORSConfig(),
		SecurityConfig:	LoadSecurityConfigs(),
	}
}
