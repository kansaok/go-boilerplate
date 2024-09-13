package cmd

import (
	"fmt"
	"os"

	"github.com/kansaok/go-boilerplate/internal/config"
	"github.com/kansaok/go-boilerplate/internal/db"
	"github.com/spf13/cobra"
)

// The rootCmd is the main command that runs the subcommands.
var rootCmd = &cobra.Command{
	Use:   "go-boilerplate",
	Short: "go-boilerplate is a simple CLI tool for database migrations",
	Long:  "go-boilerplate is a tool for managing database migrations and other operations.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
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

func init() {
	rootCmd.PersistentFlags().StringP("func", "f", "", "Run a specific seeder function")
}
