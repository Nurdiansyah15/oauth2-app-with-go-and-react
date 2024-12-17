package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string `json:"id" gorm:"primary_key" type:"uuid"`
	Username  string `json:"username" gorm:"unique"`
	Password  string `json:"-"` // Password tidak akan di-serialize ke JSON
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (user *User) BeforeCreate(db *gorm.DB) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	return nil
}

func (user *User) BeforeUpdate(db *gorm.DB) error {
	user.UpdatedAt = time.Now()
	return nil
}
