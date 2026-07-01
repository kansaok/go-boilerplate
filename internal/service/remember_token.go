package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type RememberTokenService struct {
	db *gorm.DB
}

type RememberToken struct {
	UserID    int64     `gorm:"primaryKey"`
	Token     string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}


func NewRememberTokenService(db *gorm.DB) *RememberTokenService {
	return &RememberTokenService{db: db}
}

func (s *RememberTokenService) CreateRememberToken(userID int64, token string) error {
	rememberToken := RememberToken{
		UserID:    userID,
		Token:     token,
		CreatedAt: time.Now(),
	}

	// Gunakan GORM untuk menyimpan data
	if err := s.db.Create(&rememberToken).Error; err != nil {
		return fmt.Errorf("failed to create remember token: %w", err)
	}
	return nil
}


func (s *RememberTokenService) GetRememberToken(userID int64) (string, error) {
	var rememberToken RememberToken
	// Gunakan GORM untuk mengambil data
	if err := s.db.Where("user_id = ?", userID).First(&rememberToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", fmt.Errorf("failed to get remember token: %w", err)
	}
	return rememberToken.Token, nil
}


func (s *RememberTokenService) DeleteRememberToken(userID int64) error {
	// Gunakan GORM untuk menghapus data
	if err := s.db.Where("user_id = ?", userID).Delete(&RememberToken{}).Error; err != nil {
		return fmt.Errorf("failed to delete remember token: %w", err)
	}
	return nil
}

