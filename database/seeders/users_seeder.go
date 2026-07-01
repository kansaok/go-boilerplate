package seeders

import (
	"log"
	"math/rand"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/kansaok/go-boilerplate/internal/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func generatePassword() string {
	// Generate a random password with at least 1 uppercase, 1 lowercase, 1 digit, and 1 special character
	const (
		upper       = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lower       = "abcdefghijklmnopqrstuvwxyz"
		digits      = "0123456789"
		specials    = "!@#$%^&*()-_+=<>?[]{}|"
		allChars    = upper + lower + digits + specials
	)
	rand.Seed(time.Now().UnixNano())
	length := 12
	password := make([]byte, length)
	password[0] = upper[rand.Intn(len(upper))]
	password[1] = lower[rand.Intn(len(lower))]
	password[2] = digits[rand.Intn(len(digits))]
	password[3] = specials[rand.Intn(len(specials))]
	for i := 4; i < length; i++ {
		password[i] = allChars[rand.Intn(len(allChars))]
	}

	// return string(password)
	return string("2345678232!Aa")
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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

	// Define optional values for title and gender
	titles := []string{"Mr", "Mrs", "Miss"}
	genders := []string{"L", "P"}

	data := make([][]interface{}, 20)
	for i := range data {
		email := faker.Email()
		password := generatePassword()
		hashedPassword, _ := hashPassword(password)
		address := faker.GetRealAddress()
		data[i] = []interface{}{
			email,
			hashedPassword,
			titles[rand.Intn(len(titles))],
			faker.FirstName(),
			faker.LastName(),
			faker.LastName(),
			genders[rand.Intn(len(genders))],
			faker.Date(),
			address.Address,
			faker.Phonenumber(),
			"1",
		}
	}

	if err := util.SeedTable(db, tableName, columns, data); err != nil {
		log.Fatalf("Failed to seed "+tableName+": %v", err)
	}

	log.Println(tableName + " seeded successfully!")
}
