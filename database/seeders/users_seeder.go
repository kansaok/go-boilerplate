package seeders

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	"github.com/go-faker/faker/v4"
	"github.com/kansaok/go-boilerplate/internal/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func secureRandomInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

func generatePassword() string {
	const (
		upper    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lower    = "abcdefghijklmnopqrstuvwxyz"
		digits   = "0123456789"
		specials = "!@#$%^&*()-_+=<>?[]{}|"
		allChars = upper + lower + digits + specials
	)
	length := 16
	password := make([]byte, length)
	password[0] = upper[secureRandomInt(len(upper))]
	password[1] = lower[secureRandomInt(len(lower))]
	password[2] = digits[secureRandomInt(len(digits))]
	password[3] = specials[secureRandomInt(len(specials))]
	for i := 4; i < length; i++ {
		password[i] = allChars[secureRandomInt(len(allChars))]
	}

	return string(password)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func UsersSeeder(db *gorm.DB) {
	var tableName = "users"
	columns := []string{
		"email",
		"password",
		"title",
		"first_name",
		"middle_name",
		"last_name",
		"gender",
		"bod",
		"pob",
		"phone_number",
		"created_by",
	}

	titles := []string{"Mr", "Mrs", "Miss"}
	genders := []string{"L", "P"}

	data := make([][]interface{}, 20)
	for i := range data {
		email := faker.Email()
		password := generatePassword()
		hashedPassword, err := hashPassword(password)
		if err != nil {
			log.Printf("Warning: failed to hash password for seed user %d: %v", i, err)
		}
		address := faker.GetRealAddress()
		data[i] = []interface{}{
			email,
			hashedPassword,
			titles[secureRandomInt(len(titles))],
			faker.FirstName(),
			faker.LastName(),
			faker.LastName(),
			genders[secureRandomInt(len(genders))],
			faker.Date(),
			address.Address,
			faker.Phonenumber(),
			"system",
		}
	}

	if err := util.SeedTable(db, tableName, columns, data); err != nil {
		fmt.Printf("Failed to seed %s: %v\n", tableName, err)
	}

	fmt.Printf("%s seeded successfully!\n", tableName)
}
