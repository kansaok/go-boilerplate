package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Config contains all application configurations.
type Config struct {
    DBDriver	string
    DBHost     	string
    DBPort     	string
    DBUser     	string
    DBPassword 	string
    DBName     	string
    DBSSLMode  	string
    DBFile  	string
	BucketName string
}

// LoadConfig loads configurations from environment variables.
func LoadDatabaseConfig() *Config {
	// Load environment variables.
    err := godotenv.Load()
	if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

	// Set Gin to the mode based on environment variables.
	ginMode := getEnv("GIN_MODE","debug")
	if ginMode == "" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(ginMode)
	}

    return &Config{
        DBDriver:	getEnv("DB_CONNECTION", "postgres"),
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "user"),
        DBPassword: getEnv("DB_PASSWORD", "password"),
        DBName:     getEnv("DB_NAME", "dbname"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
        DBFile:  	getEnv("DB_FILE", "disable"),
		BucketName: getEnv("S3_BUCKET_NAME","s3bucket"),
    }
}

// getEnv retrieves the value of an environment variable with a default fallback
func getEnv(key string, defaultValue string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        return defaultValue
    }
    return value
}
