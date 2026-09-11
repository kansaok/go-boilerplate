package service

import (
	"crypto/sha256"
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

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}

func (s *RememberTokenService) CreateRememberToken(userID int64, token string) error {
	rememberToken := RememberToken{
		UserID:    userID,
		Token:     hashToken(token),
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(&rememberToken).Error; err != nil {
		return fmt.Errorf("failed to create remember token: %w", err)
	}
	return nil
}

func (s *RememberTokenService) GetRememberToken(userID int64) (string, error) {
	var rememberToken RememberToken
	if err := s.db.Where("user_id = ?", userID).First(&rememberToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil
		}
		return "", fmt.Errorf("failed to get remember token: %w", err)
	}
	return rememberToken.Token, nil
}

func (s *RememberTokenService) VerifyRememberToken(userID int64, token string) (bool, error) {
	hashed := hashToken(token)
	var rememberToken RememberToken
	if err := s.db.Where("user_id = ? AND token = ?", userID, hashed).First(&rememberToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to verify remember token: %w", err)
	}
	return true, nil
}

func (s *RememberTokenService) DeleteRememberToken(userID int64) error {
	if err := s.db.Where("user_id = ?", userID).Delete(&RememberToken{}).Error; err != nil {
		return fmt.Errorf("failed to delete remember token: %w", err)
	}
	return nil
}
