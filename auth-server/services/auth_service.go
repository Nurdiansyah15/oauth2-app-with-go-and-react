package services

import (
	"auth-server/models"
	"auth-server/repositories"
)

type AuthService interface {
	Login(string, string) (*models.User, string)
}

type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) *authService {
	return &authService{userRepository: userRepository}
}

func (s *authService) Login(username, password string) (*models.User, string) {
	user, errCode := s.userRepository.GetUserByUsername(username)
	if errCode != "" {
		return nil, errCode
	}
	if user.Password != password {
		return nil, "WRONG_PASSWORD_401"
	}
	return user, ""
}

// func Logout(c *gin.Context) {
// 	session := sessions.Default(c)
// 	session.Delete("user_id")
// 	session.Delete("username")
// 	session.Save()
// }
