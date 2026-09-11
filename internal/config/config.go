package config

import (
	"log"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

// AppConfig stores all application configurations.
type AppConfig struct {
	DatabaseConfig *Config
	JWTConfig      *JWTConfig
	CORSConfig     cors.Config
	SecurityConfig SecurityConfig
	RedisConfig    *RedisConfig
}

var (
	appConfig  *AppConfig
	configOnce sync.Once
)

// ResetConfigTestOnly mereset cache konfigurasi agar unit test bisa memuat
// ulang environment. Hanya boleh dipanggil dari test.
func ResetConfigTestOnly() {
	configOnce = sync.Once{}
}

// LoadConfig loads all application configurations (cached after first call).
func LoadConfig() *AppConfig {
	configOnce.Do(func() {
		// Environment variables may come from the process (Docker env_file) or a local .env file.
		// The .env file is optional: ignore load errors and rely on real environment variables.
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: .env file not found, relying on environment variables")
		}

		appConfig = &AppConfig{
			DatabaseConfig: LoadDatabaseConfig(),
			JWTConfig:      LoadJWTConfig(),
			CORSConfig:     CORSConfig(),
			SecurityConfig: LoadSecurityConfigs(),
			RedisConfig:    LoadRedisConfig(),
		}
	})
	return appConfig
}