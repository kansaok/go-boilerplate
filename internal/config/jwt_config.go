package config

import (
	"log"
	"os"
	"time"
)

type JWTConfig struct {
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
	SecretKey            string
}

func LoadJWTConfig() *JWTConfig {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("FATAL: JWT_SECRET_KEY environment variable is required but not set")
	}
	if len(secretKey) < 32 {
		log.Fatal("FATAL: JWT_SECRET_KEY must be at least 32 characters long")
	}

	return &JWTConfig{
		AccessTokenLifetime:  getEnvAsDuration("ACCESS_TOKEN_LIFETIME", time.Minute*15),
		RefreshTokenLifetime: getEnvAsDuration("REFRESH_TOKEN_LIFETIME", time.Hour*24*7),
		SecretKey:            secretKey,
	}
}

func getEnvAsDuration(name string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(name); exists {
		d, err := time.ParseDuration(value)
		if err == nil {
			return d
		}
	}
	return defaultValue
}
