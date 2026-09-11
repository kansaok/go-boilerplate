package config

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// Config contains all application configurations.
type Config struct {
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBFile     string
	BucketName string
}

// LoadDatabaseConfig loads configurations from environment variables.
func LoadDatabaseConfig() *Config {
	// Set Gin to the mode based on environment variables.
	ginMode := getEnv("GIN_MODE", "debug")
	switch ginMode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// Validate required fields for SQL databases
	dbDriver := getEnv("DB_CONNECTION", "postgres")
	switch dbDriver {
	case "postgres", "mysql", "sqlite":
		if dbDriver == "sqlite" {
			dbFile := getEnv("DB_FILE", "")
			if dbFile == "" {
				log.Fatal("FATAL: DB_FILE environment variable is required for SQLite")
			}
			return &Config{
				DBDriver:   dbDriver,
				DBFile:     dbFile,
				BucketName: getEnv("S3_BUCKET_NAME", ""),
			}
		}
		dbUser := getEnv("DB_USER", "")
		dbPassword := getEnv("DB_PASSWORD", "")
		dbName := getEnv("DB_NAME", "")
		if dbUser == "" || dbPassword == "" || dbName == "" {
			log.Fatal("FATAL: DB_USER, DB_PASSWORD, and DB_NAME are required for SQL databases")
		}
		return &Config{
			DBDriver:   dbDriver,
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			DBUser:     dbUser,
			DBPassword: dbPassword,
			DBName:     dbName,
			DBSSLMode:  getEnv("DB_SSLMODE", "require"),
			BucketName: getEnv("S3_BUCKET_NAME", ""),
		}
	case "mongodb":
		dbHost := getEnv("DB_HOST", "")
		dbPort := getEnv("DB_PORT", "")
		dbUser := getEnv("DB_USER", "")
		dbPassword := getEnv("DB_PASSWORD", "")
		dbName := getEnv("DB_NAME", "")
		if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
			log.Fatal("FATAL: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, and DB_NAME are required for MongoDB")
		}
		return &Config{
			DBDriver:   dbDriver,
			DBHost:     dbHost,
			DBPort:     dbPort,
			DBUser:     dbUser,
			DBPassword: dbPassword,
			DBName:     dbName,
			BucketName: getEnv("S3_BUCKET_NAME", ""),
		}
	default:
		log.Fatalf("FATAL: unsupported database driver: %s", dbDriver)
		return nil
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
