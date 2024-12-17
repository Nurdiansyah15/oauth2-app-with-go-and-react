package repositories

import (
	"auth-server/database"
	"auth-server/models"
)

func CreateAuthCode(authCode models.AuthCode) error {
	return database.DB.Create(&authCode).Error
}

func FindAuthCodeByCode(code string) (*models.AuthCode, error) {
	var authCode models.AuthCode
	result := database.DB.Where("code = ?", code).First(&authCode)
	return &authCode, result.Error
}

func DeleteAuthCode(code string) error {
	return database.DB.Where("code = ?", code).Delete(&models.AuthCode{}).Error
}
