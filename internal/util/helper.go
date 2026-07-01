package util

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kansaok/go-boilerplate/pkg/logger"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

func UpdateSeederMap(functionName string) error {
	seederFilePath := "internal/db/seeder.go"

	// Read the existing seeder.go file
	data, err := os.ReadFile(seederFilePath)
	if err != nil {
		return fmt.Errorf("could not read seeder.go: %v", err)
	}

	// Convert data to a string
	content := string(data)

	// Find the position of the seederMap variable
	mapPosition := strings.Index(content, "var seederMap = map[string]func(*gorm.DB)")
	if mapPosition == -1 {
		return fmt.Errorf("could not find seederMap declaration")
	}

	// Find the position of the closing curly brace of the map
	closeBracePosition := strings.Index(content[mapPosition:], "}")
	if closeBracePosition == -1 {
		return fmt.Errorf("could not find the end of the seederMap")
	}
	closeBracePosition += mapPosition

	// Create the new entry for the seeder
	newSeederEntry := fmt.Sprintf(`"%s": seeders.%s,`, functionName, functionName)

	// Insert the new entry before the closing brace
	newContent := content[:closeBracePosition] + "\n    " + newSeederEntry + content[closeBracePosition:]

	// Write the updated content back to seeder.go
	if err := os.WriteFile(seederFilePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("could not write to seeder.go: %v", err)
	}

	return nil
}

func SeedTable(db *gorm.DB, tableName string, columns []string, data [][]interface{}) error {
	// Mulai transaksi untuk menjaga konsistensi data
	tx := db.Begin() // Mulai transaksi
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Hapus data yang sudah ada di tabel
	if err := tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Buat query insert
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "?" // GORM menggunakan tanda tanya untuk placeholders
	}
	insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	// Insert data ke dalam tabel
	for _, row := range data {
		if err := tx.Exec(insertQuery, row...).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func ToPascalCase(input string) string {
	words := strings.Split(input, "_")
	titleCaser := cases.Title(language.English)
	for i := range words {
		words[i] = titleCaser.String(words[i])
	}
	return strings.Join(words, "")
}

func HandleTransaction(ctx context.Context, tx *sql.Tx) {
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logger.Critical(ctx, fmt.Sprintf("transaksi dibatalkan: %v", r))
		}
	}()
}

func ParseDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02" // Expected format
	return time.Parse(layout, dateStr)
}
