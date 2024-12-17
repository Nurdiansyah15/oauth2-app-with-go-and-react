package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Session struct {
	ID           string    `json:"id" gorm:"primary_key" type:"uuid"`
	UserID       string    `json:"user_id"`
	ClientID     string    `json:"client_id"`
	BrowserInfo  string    `json:"browser_info"`
	IPAddress    string    `json:"ip_address"`
	SessionToken string    `json:"session_token"`
	LastActivity time.Time `json:"last_activity"`
	ExpiredAt    time.Time `json:"expired_at"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (session *Session) BeforeCreate(db *gorm.DB) error {
	session.ID = uuid.New().String()
	session.CreatedAt = time.Now()
	return nil
}

func (session *Session) BeforeUpdate(db *gorm.DB) error {
	session.UpdatedAt = time.Now()
	return nil
}
