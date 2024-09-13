package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/kansaok/go-boilerplate/internal/db"
	"github.com/kansaok/go-boilerplate/internal/util"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var funcName string

var makeSeederCmd = &cobra.Command{
	Use:   "make:seeder",
	Short: "Generate a new seeder file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		seederName := args[0]
		// Convert seeder name to PascalCase for function name
		functionName := util.ToPascalCase(seederName)
		fileName := fmt.Sprintf("database/seeders/%s.go", seederName)

		// Check if the file already exists
		if _, err := os.Stat(fileName); err == nil {
			log.Fatalf("Seeder %s already exists!", fileName)
		}

		// Create seeder file
		content := fmt.Sprintf(`package seeders

import (
	"log"

	"github.com/kansaok/go-boilerplate/internal/util"
	"gorm.io/gorm"
)

func %s(db *gorm.DB) {
	var tableName = "table_name"
	columns := []string{"list", "column", "here"}
	data := [][]interface{}{
		{"data", "or value", "here"},
	}

	if err := util.SeedTable(db, tableName, columns, data); err != nil {
		log.Fatalf("Failed to seed "+tableName+": %%v", err)
	}

	log.Println(tableName + " seeded successfully!")
}
`, functionName)

		// Create the file with the template
		if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
			log.Fatalf("Failed to create seeder file: %v", err)
		}

		// Update the seeder map in internal/db/seeder.go
		if err := util.UpdateSeederMap(functionName); err != nil {
			log.Fatalf("Failed to update seeder map: %v", err)
		}

		log.Printf("Seeder created: %s", fileName)
	},
}

var seedCmd = &cobra.Command{
    Use:   "db:seed",
    Short: "Run database seeders",
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
			// Cek apakah ada flag --func
			if funcName != "" {
				if err := db.RunSpecificSeeder(sqlDB, funcName); err != nil {
					log.Fatalf("Failed to run specific seeder: %v", err)
				}
			} else {
				if err := db.RunAllSeeders(sqlDB); err != nil {
					log.Fatalf("Failed to run seeders: %v", err)
				}
			}
			log.Println("Database seeding completed successfully!")
		} else {
			log.Fatalf("Expected *gorm.DB but got %T", dbConn)
		}
    },
}

func init() {
    rootCmd.AddCommand(makeSeederCmd)
	seedCmd.Flags().StringVar(&funcName, "func", "", "Run a specific seeder function")
    rootCmd.AddCommand(seedCmd)
}
