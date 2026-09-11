package util

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/kansaok/go-boilerplate/pkg/logger"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

var validIdentifier = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func sanitizeIdentifier(name string) error {
	if !validIdentifier.MatchString(name) {
		return fmt.Errorf("invalid identifier: %s", name)
	}
	return nil
}

func UpdateSeederMap(functionName string) error {
	if err := sanitizeIdentifier(functionName); err != nil {
		return fmt.Errorf("invalid function name: %v", err)
	}

	seederFilePath := "internal/db/seeder.go"

	data, err := os.ReadFile(seederFilePath)
	if err != nil {
		return fmt.Errorf("could not read seeder.go: %v", err)
	}

	content := string(data)

	mapPosition := strings.Index(content, "var seederMap = map[string]func(*gorm.DB)")
	if mapPosition == -1 {
		return fmt.Errorf("could not find seederMap declaration")
	}

	closeBracePosition := strings.Index(content[mapPosition:], "}")
	if closeBracePosition == -1 {
		return fmt.Errorf("could not find the end of the seederMap")
	}
	closeBracePosition += mapPosition

	newSeederEntry := fmt.Sprintf(`"%s": seeders.%s,`, functionName, functionName)

	newContent := content[:closeBracePosition] + "\n    " + newSeederEntry + content[closeBracePosition:]

	if err := os.WriteFile(seederFilePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("could not write to seeder.go: %v", err)
	}

	return nil
}

func SeedTable(db *gorm.DB, tableName string, columns []string, data [][]interface{}) error {
	if err := sanitizeIdentifier(tableName); err != nil {
		return fmt.Errorf("invalid table name: %v", err)
	}
	for _, col := range columns {
		if err := sanitizeIdentifier(col); err != nil {
			return fmt.Errorf("invalid column name: %v", err)
		}
	}

	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Exec(fmt.Sprintf("TRUNCATE TABLE %s", tableName)).Error; err != nil {
		tx.Rollback()
		return err
	}

	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	for _, row := range data {
		if err := tx.Exec(insertQuery, row...).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

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
	layout := "2006-01-02"
	return time.Parse(layout, dateStr)
}
