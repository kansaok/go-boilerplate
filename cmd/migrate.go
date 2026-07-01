package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/kansaok/go-boilerplate/internal/db"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var makeMigrationCmd = &cobra.Command{
	Use:   "make:migration",
	Short: "Generate a new migration file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		migrationName := args[0]
		timestamp := time.Now().Format("20060102150405")

		upFileName := fmt.Sprintf("database/migrations/%s_%s.sql", timestamp, migrationName)
		downFileName := fmt.Sprintf("database/migrations/%s_%s_rollback.sql", timestamp, migrationName)

		if err := db.CreateMigrationFile(upFileName); err != nil {
			log.Fatalf("Failed to create migration file: %v", err)
		}
		if err := db.CreateMigrationFile(downFileName); err != nil {
			log.Fatalf("Failed to create rollback migration file: %v", err)
		}

		log.Printf("Migration created: %s and %s", upFileName, downFileName)
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		config := loadConfig()
		dbConn, err := db.ConnectDB(config)
		if err != nil {
			log.Fatalf("Failed to connect to the database: %v", err)
		}

		// Type assertion for *gorm.DB
		if sqlDB, ok := dbConn.(*gorm.DB); ok {
			sqlConn, err := sqlDB.DB()
			if err != nil {
				log.Fatalf("Failed to get database connection: %v", err)
			}
			defer sqlConn.Close()
			if err := db.RunMigrations(sqlDB); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
		} else {
			log.Fatalf("Expected *gorm.DB but got %T", dbConn)
		}

		log.Println("Migration applied successfully!")
	},
}

var rollbackCmd = &cobra.Command{
	Use:   "migrate:rollback",
	Short: "Rollback the last database migration",
	Run: func(cmd *cobra.Command, args []string) {
		config := loadConfig()
		dbConn, err := db.ConnectDB(config)
		if err != nil {
			log.Fatalf("Failed to connect to the database: %v", err)
		}

		// Type assertion for *gorm.DB
		if sqlDB, ok := dbConn.(*gorm.DB); ok {
			sqlConn, err := sqlDB.DB()
			if err != nil {
				log.Fatalf("Failed to get database connection: %v", err)
			}
			defer sqlConn.Close()
			if err := db.RollbackLastMigration(sqlDB); err != nil {
				log.Fatalf("Rollback failed: %v", err)
			}
		} else {
			log.Fatalf("Expected *gorm.DB but got %T", dbConn)
		}

		log.Println("Migration rollback successful!")
	},
}

var statusCmd = &cobra.Command{
	Use:   "migrate:status",
	Short: "Show migration status",
	Run: func(cmd *cobra.Command, args []string) {
		config := loadConfig()
		dbConn, err := db.ConnectDB(config)
		if err != nil {
			log.Fatalf("Failed to connect to the database: %v", err)
		}

		// Type assertion for *gorm.DB
		if sqlDB, ok := dbConn.(*gorm.DB); ok {
			sqlConn, err := sqlDB.DB()
			if err != nil {
				log.Fatalf("Failed to get database connection: %v", err)
			}
			defer sqlConn.Close()
			if err := db.ShowMigrationStatus(sqlDB); err != nil {
				log.Fatalf("Status check failed: %v", err)
			}
		} else {
			log.Fatalf("Expected *gorm.DB but got %T", dbConn)
		}
	},
}

var migrateFreshCmd = &cobra.Command{
    Use:   "migrate:fresh",
    Short: "Drop all tables and re-run all migrations",
    Run: func(cmd *cobra.Command, args []string) {
		config := loadConfig()
		dbConn, err := db.ConnectDB(config)
		if err != nil {
			log.Fatalf("Failed to connect to the database: %v", err)
		}

		// Type assertion for *gorm.DB
		if sqlDB, ok := dbConn.(*gorm.DB); ok {
			sqlConn, err := sqlDB.DB()
			if err != nil {
				log.Fatalf("Failed to get database connection: %v", err)
			}
			defer sqlConn.Close()
			if err := db.FreshMigrate(sqlDB, config.DBName); err != nil {
				log.Fatalf("Status check failed: %v", err)
			}
		} else {
			log.Fatalf("Expected *gorm.DB but got %T", dbConn)
		}
    },
}

func init() {
	rootCmd.AddCommand(makeMigrationCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(rollbackCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(migrateFreshCmd)
}
