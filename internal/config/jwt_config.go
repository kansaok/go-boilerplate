package config

import (
	"os"
	"time"
)

type JWTConfig struct {
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
	SecretKey            string
}

func LoadJWTConfig() *JWTConfig {
	return &JWTConfig{
		AccessTokenLifetime:  getEnvAsDuration("ACCESS_TOKEN_LIFETIME", time.Minute*15),
		RefreshTokenLifetime: getEnvAsDuration("REFRESH_TOKEN_LIFETIME", time.Hour*24*7),
		SecretKey:            getEnv("JWT_SECRET_KEY", "supersecretkey"),
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
