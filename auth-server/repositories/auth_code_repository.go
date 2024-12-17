package repositories

import (
	"auth-server/models"

	"gorm.io/gorm"
)

type AuthCodeRepository interface {
	CreateAuthCode(authCode *models.AuthCode) string
	GetAuthCodeByCode(code string) (*models.AuthCode, string)
	DeleteAuthCode(code string) string
}

type authCodeRepository struct {
	db *gorm.DB
}

func NewAuthCodeRepository(db *gorm.DB) AuthCodeRepository {
	return &authCodeRepository{db: db}
}

func (r *authCodeRepository) CreateAuthCode(authCode *models.AuthCode) string {
	if err := r.db.Create(authCode).Error; err != nil {
		return "DATABASE_ERROR_500"
	}
	return ""
}

func (r *authCodeRepository) GetAuthCodeByCode(code string) (*models.AuthCode, string) {
	var authCode models.AuthCode

	if err := r.db.Where("code = ?", code).First(&authCode).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "AUTH_CODE_NOT_FOUND_404"
		}
		return nil, "DATABASE_ERROR_500"
	}

	return &authCode, ""
}

func (r *authCodeRepository) DeleteAuthCode(code string) string {
	result := r.db.Where("code = ?", code).Delete(&models.AuthCode{})

	if result.Error != nil {
		return "DATABASE_ERROR_500"
	}

	if result.RowsAffected == 0 {
		return "AUTH_CODE_NOT_FOUND_404"
	}

	return ""
}
