package repositories

import (
	"auth-server/models"
	"errors"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserById(id string) (*models.User, string)
	GetUserByUsername(username string) (*models.User, string)
	GetAllUsers() ([]*models.User, string)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUserById(id string) (*models.User, string) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "USER_NOT_FOUND_404"
		}
		return nil, "DATABASE_ERROR_500"
	}
	return &user, ""
}

func (r *userRepository) GetUserByUsername(username string) (*models.User, string) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "USER_NOT_FOUND_404"
		}
		return nil, "DATABASE_ERROR_500"
	}
	return &user, ""
}

func (r *userRepository) GetAllUsers() ([]*models.User, string) {
	var users []*models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, "DATABASE_ERROR_500"
	}
	return users, ""
}
