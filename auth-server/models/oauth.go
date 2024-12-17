package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthCode struct {
	ID        string    `json:"id" gorm:"primary_key" type:"uuid"`
	Code      string    `json:"code"`
	ClientID  string    `json:"client_id"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (authCode *AuthCode) BeforeCreate(db *gorm.DB) error {
	authCode.ID = uuid.New().String()
	authCode.CreatedAt = time.Now()
	return nil
}

func (authCode *AuthCode) BeforeUpdate(db *gorm.DB) error {
	authCode.UpdatedAt = time.Now()
	return nil
}
