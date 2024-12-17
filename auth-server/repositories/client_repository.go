package repositories

import (
	"auth-server/models"

	"gorm.io/gorm"
)

type ClientRepository interface {
	GetClientByID(clientID string) (*models.Client, string)
}

type clientRepository struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) GetClientByID(clientID string) (*models.Client, string) {
	var client models.Client

	if err := r.db.Where("id = ?", clientID).First(&client).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "CLIENT_NOT_FOUND_404"
		}
		return nil, "DATABASE_ERROR_500"
	}

	return &client, ""
}
