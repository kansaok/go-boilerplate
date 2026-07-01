package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kansaok/go-boilerplate/internal/config"
	"github.com/kansaok/go-boilerplate/internal/db"
	"github.com/kansaok/go-boilerplate/internal/routes"
	"github.com/kansaok/go-boilerplate/internal/util/validators"
	"github.com/kansaok/go-boilerplate/pkg/logger"
	"github.com/kansaok/go-boilerplate/pkg/telemetry"

	"github.com/kansaok/go-boilerplate/cmd"
)

func main() {
	cmdMap := map[string]func(){
		"make:migration":  cmd.Execute,
		"make:seeder":  cmd.Execute,
		"migrate":       cmd.Execute,
		"migrate:fresh": cmd.Execute,
		"migrate:rollback": cmd.Execute,
		"migrate:status": cmd.Execute,
		"db:seed": cmd.Execute,
		"-h": cmd.Execute,
		"--help": cmd.Execute,
		"--func": cmd.Execute,
	}

	// Cek argumen command line
	if len(os.Args) > 1 {
		command := strings.ToLower(os.Args[1])
		if execFunc, exists := cmdMap[command]; exists {
			execFunc()
			return
		}
	}

	runServer()
}

func runServer() {
	logger.Init()
	validators.LoadValidatorConfig()
	// err := http.ListenAndServeTLS(":443", "cert.pem", "key.pem", router)
    // if err != nil {
    //     log.Fatal("ListenAndServeTLS: ", err)
    // }

    // Disable directory listing
    fs := http.FileServer(http.Dir("./web/assets"))
    http.Handle("/assets/", http.StripPrefix("/assets/", fs))

    // Convert to database configuration
    dbConfig := loadConfig()

    // Connect to the database
    _, err := db.ConnectDB(dbConfig)
    if err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }

    // Initialize OpenTelemetry
    // cleanup := telemetry.InitOpenTelemetry("myapp")
    // defer cleanup()

	// rememberTokenService := service.NewRememberTokenService(db)

	// Initialize Prometheus
    telemetry.InitPrometheus()

    // Define routes
    r := routes.SetupRoutes()

    // Run the application
    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func loadConfig() *db.Config {
	cfg := config.LoadConfig()
	return &db.Config{
		DBHost:       cfg.DatabaseConfig.DBHost,
		DBPort:       cfg.DatabaseConfig.DBPort,
		DBUser:       cfg.DatabaseConfig.DBUser,
		DBPassword:   cfg.DatabaseConfig.DBPassword,
		DBName:       cfg.DatabaseConfig.DBName,
		DBSSLMode:    cfg.DatabaseConfig.DBSSLMode,
		DBConnection: cfg.DatabaseConfig.DBDriver,
		DBFile:       cfg.DatabaseConfig.DBFile,
	}
}
