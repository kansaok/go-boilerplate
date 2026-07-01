package user

import (
	"time"
)

type User struct {
    ID          int64     `gorm:"primaryKey;autoIncrement;index;unique"`
    Email       string    `gorm:"size:100;not null;unique"`
    Password    string    `gorm:"size:200;not null"`
    Title       string    `gorm:"type:enum('Mr', 'Mrs', 'Miss');not null"`
    FirstName   string    `gorm:"size:50;not null"`
    MiddleName  *string   `gorm:"size:50"`
    LastName    *string   `gorm:"size:50"`
    Gender      string    `gorm:"type:enum('L', 'P');not null"`
    Bod         time.Time `gorm:"not null"`
    Pob         string    `gorm:"size:100;not null"`
    PhoneNumber string    `gorm:"size:20;not null"`
    Photo       *string   `gorm:"size:200"`
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    CreatedBy   string    `gorm:"size:100;not null"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
    UpdatedBy   *string   `gorm:"size:100"`
}
