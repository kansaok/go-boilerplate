package config

import (
	"log"
	"os"
	"strconv"
)

// RedisConfig stores optional Redis connection settings used for
// cross-instance rate limiting (e.g. per-account login/register throttle).
type RedisConfig struct {
	Enabled  bool
	Host     string
	Port     string
	Password string
	DB       int
}

// LoadRedisConfig reads Redis settings from environment variables.
// Redis is optional: when REDIS_HOST is empty the app falls back to an
// in-memory limiter (suitable for single-instance deployments).
func LoadRedisConfig() *RedisConfig {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		return &RedisConfig{Enabled: false}
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	db := 0
	if raw := os.Getenv("REDIS_DB"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			log.Printf("WARNING: invalid REDIS_DB %q, defaulting to 0", raw)
		} else {
			db = parsed
		}
	}

	return &RedisConfig{
		Enabled:  true,
		Host:     host,
		Port:     port,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	}
}