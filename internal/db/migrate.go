package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

var migrationsDir = "database/migrations"

// CreateMigrationsTable creates the migrations table if it doesn't already exist.
func CreateMigrationsTable(db *gorm.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS remember_tokens (
		user_id BIGINT NOT NULL,
		token TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		PRIMARY KEY (user_id)
	);
	CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		applied_at TIMESTAMP NOT NULL
	)`
	if err := db.Exec(query).Error; err != nil {
		return fmt.Errorf("failed to create root table: %w", err)
	}

	return nil
}

func CreateMigrationFile(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("-- SQL migration script here")
	if err != nil {
		return err
	}

	return nil
}

// RunMigrations executes all .sql files that have not been applied yet
func RunMigrations(db *gorm.DB) error {
	if err := CreateMigrationsTable(db); err != nil {
		return err
	}

	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	// Use os.ReadDir to read the list of files in the directory.
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") && !strings.Contains(file.Name(), "_rollback.sql") {
			migrationName := file.Name()
			if _, applied := appliedMigrations[migrationName]; !applied {
				log.Printf("Applying migration: %s", migrationName)
				if err := applySQLFile(db, filepath.Join(migrationsDir, migrationName)); err != nil {
					return fmt.Errorf("failed to apply migration %s: %w", migrationName, err)
				}
				if err := recordMigration(db, migrationName); err != nil {
					return fmt.Errorf("failed to record migration %s: %w", migrationName, err)
				}
			}
		}
	}

	return nil
}

// RollbackLastMigration reverts the last applied migration.
func RollbackLastMigration(db *gorm.DB) error {
	if err := CreateMigrationsTable(db); err != nil {
		return err
	}

	// Retrieve the name of the last applied migration.
	lastMigration, err := getLastAppliedMigration(db)
	if err != nil {
		return err
	}
	if lastMigration == "" {
		log.Println("No migrations to rollback")
		return nil
	}

	// Rename the migration file to the rollback file.
	rollbackFile := strings.Replace(lastMigration, ".sql", "_rollback.sql", 1)
	rollbackFilePath := filepath.Join(migrationsDir, rollbackFile)

	// Ensure that the rollback file exists.
	if _, err := os.Stat(rollbackFilePath); os.IsNotExist(err) {
		return fmt.Errorf("rollback file not found for migration %s", lastMigration)
	}

	// Execute the rollback SQL file.
	log.Printf("Rolling back migration: %s", lastMigration)
	if err := applySQLFile(db, rollbackFilePath); err != nil {
		return fmt.Errorf("failed to rollback migration %s: %w", lastMigration, err)
	}

	// Remove the migration record from the migrations table.
	if err := deleteMigrationRecord(db, lastMigration); err != nil {
		return fmt.Errorf("failed to delete migration record %s: %w", lastMigration, err)
	}

	return nil
}

// applySQLFile runs the SQL file.
func applySQLFile(db *gorm.DB, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read SQL file %s: %w", filePath, err)
	}

	// Use db.Exec and handle the error from execution.
	err = db.Exec(string(content)).Error
	if err != nil {
		return fmt.Errorf("failed to execute SQL file %s: %w", filePath, err)
	}

	return nil
}


// ShowMigrationStatus displays the status of all migrations.
func ShowMigrationStatus(db *gorm.DB) error {
	if err := CreateMigrationsTable(db); err != nil {
		return err
	}

	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return err
	}

	// Use os.ReadDir to read the list of files in the directory.
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	log.Println("Migration status:")
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") && !strings.Contains(file.Name(), "_rollback.sql") {
			if _, applied := appliedMigrations[file.Name()]; applied {
				log.Printf("[X] %s", file.Name())
			} else {
				log.Printf("[ ] %s", file.Name())
			}
		}
	}

	return nil
}

// getAppliedMigrations retrieves the list of migrations that have been applied.
func getAppliedMigrations(db *gorm.DB) (map[string]time.Time, error) {
	rows, err := db.Raw("SELECT name, applied_at FROM migrations").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appliedMigrations := make(map[string]time.Time)
	for rows.Next() {
		var name string
		var appliedAt time.Time
		if err := rows.Scan(&name, &appliedAt); err != nil {
			return nil, err
		}
		appliedMigrations[name] = appliedAt
	}
	return appliedMigrations, nil
}

// getLastAppliedMigration gets the name of the last applied migration.
func getLastAppliedMigration(db *gorm.DB) (string, error) {
	var name string
	err := db.Raw("SELECT name FROM migrations ORDER BY applied_at DESC LIMIT 1").Scan(&name).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", err
	}

	return name, nil
}

// recordMigration logs the migrations that have been applied.
func recordMigration(db *gorm.DB, name string) error {
	err := db.Exec("INSERT INTO migrations (name, applied_at) VALUES (?, ?)", name, time.Now()).Error
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}
	return nil
}

// deleteMigrationRecord removes the migration record.
func deleteMigrationRecord(db *gorm.DB, name string) error {
	err := db.Exec("DELETE FROM migrations WHERE name = ?", name).Error
	if err != nil {
		return fmt.Errorf("failed to delete migration record: %w", err)
	}
	return nil
}

// FreshMigrate drops all tables and runs all migrations from the beginning.
func FreshMigrate(db *gorm.DB,dbName string) error {
    // Drop all tables
    if err := dropAllTables(db); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    RunMigrations(db)

	return nil
}

func dropAllTables(db *gorm.DB) error {
    query := `
    DO $$
    DECLARE
        r RECORD;
    BEGIN
        FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
            EXECUTE 'DROP TABLE IF EXISTS ' || r.tablename || ' CASCADE';
        END LOOP;
    END $$;`

    err := db.Exec(query)
    if err != nil {
        return fmt.Errorf("failed to drop all tables: %v", err)
    }

    fmt.Println("All tables dropped successfully!")
    return nil
}
