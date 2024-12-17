package services

import (
	"auth-server/models"
	"auth-server/repositories"
)

type UserService interface {
	GetUserById(string) (*models.User, string)
	GetAllUsers() ([]*models.User, string)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (s *userService) GetUserById(id string) (*models.User, string) {
	return s.userRepository.GetUserById(id)
}

func (s *userService) GetAllUsers() ([]*models.User, string) {
	return s.userRepository.GetAllUsers()
}
