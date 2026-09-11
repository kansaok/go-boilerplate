package auth

import (
	"context"
	"fmt"

	"github.com/kansaok/go-boilerplate/internal/db"
	usr "github.com/kansaok/go-boilerplate/internal/modules/user"
	"gorm.io/gorm"
)

func CreateUser(ctx context.Context, user usr.User) (map[string]interface{}, error) {
	dbConn, ok := db.GetDB().(*gorm.DB)
    if !ok {
        return nil, fmt.Errorf("gagal mengonversi koneksi ke *gorm.DB")
    }

    err := dbConn.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(&user).Error; err != nil {
            return fmt.Errorf("gagal menyimpan pengguna: %w", err)
        }
        return nil
    })

    if err != nil {
        return nil, err
    }

    var result map[string]interface{}
    err = dbConn.Model(&user).
        Select("id, first_name, email, title, gender, bod, pob, phone_number").
        Where("id = ?", user.ID).
        Scan(&result).Error
    if err != nil {
        return nil, fmt.Errorf("gagal mengambil data pengguna: %w", err)
    }

    return result, nil
}

func GetUserByEmail(email string) (*usr.User, error) {
    dbConn, ok := db.GetDB().(*gorm.DB)
    if !ok {
        return nil, fmt.Errorf("gagal mengonversi koneksi ke *gorm.DB")
    }

    var user usr.User
    if err := dbConn.Where("email = ?", email).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, fmt.Errorf("gagal mengambil pengguna: %w", err)
    }

    return &user, nil
}
