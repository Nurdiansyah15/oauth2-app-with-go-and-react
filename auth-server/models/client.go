package models

import (
	"time"

	"gorm.io/gorm"
)

type Client struct {
	ID          string `json:"id" gorm:"primary_key"`
	Secret      string `json:"secret"`
	RedirectURI string `json:"redirect_uri"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (client *Client) BeforeCreate(db *gorm.DB) error {
	// client.ID = uuid.New().String()
	client.CreatedAt = time.Now()
	return nil
}

func (client *Client) BeforeUpdate(db *gorm.DB) error {
	client.UpdatedAt = time.Now()
	return nil
}
