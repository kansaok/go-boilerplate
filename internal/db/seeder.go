package db

import (
	"fmt"

	"github.com/kansaok/go-boilerplate/database/seeders"
	"gorm.io/gorm"
)

var seederMap = map[string]func(*gorm.DB){
    "UsersSeeder":   seeders.UsersSeeder,
}

// RunSpecificSeeder runs a seeder based on the function name.
func RunSpecificSeeder(db *gorm.DB, funcName string) error {
    seeder, exists := seederMap[funcName]
    if !exists {
        return fmt.Errorf("seeder %s not found", funcName)
    }

    seeder(db)
    return nil
}

// RunAllSeeders runs all registered seeders.
func RunAllSeeders(db *gorm.DB) error {
    for _, seeder := range seederMap {
        seeder(db)
    }

    return nil
}
