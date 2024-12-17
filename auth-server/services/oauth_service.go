package services

import (
	"auth-server/models"
	"auth-server/repositories"
	"auth-server/utils"
	"time"
)

type OAuthService interface {
	ValidateClient(clientID, redirectURI string) (*models.Client, string)
	ValidateSession(sessionToken string) (*models.Session, string)
	ValidateAuthCode(code, clientID string) (*models.AuthCode, string)
	GenerateAuthCodeWithSession(userID, clientID string) (*models.AuthCode, string)
	ExchangeAuthCodeForToken(code, clientID, clientSecret string) (map[string]interface{}, string)
	ValidateAccessToken(tokenString string) (map[string]string, string)
}

type oauthService struct {
	clientRepo   repositories.ClientRepository
	sessionRepo  repositories.SessionRepository
	authCodeRepo repositories.AuthCodeRepository
	jwtService   JWTService
}

func NewOAuthService(
	clientRepo repositories.ClientRepository,
	sessionRepo repositories.SessionRepository,
	authCodeRepo repositories.AuthCodeRepository,
	jwtService JWTService,
) OAuthService {
	return &oauthService{
		clientRepo:   clientRepo,
		sessionRepo:  sessionRepo,
		authCodeRepo: authCodeRepo,
		jwtService:   jwtService,
	}
}

func (s *oauthService) ValidateClient(clientID, redirectURI string) (*models.Client, string) {
	client, errCode := s.clientRepo.GetClientByID(clientID)
	if errCode != "" {
		return nil, "INVALID_CLIENT"
	}

	if client.RedirectURI != redirectURI {
		return nil, "INVALID_REDIRECT_URI"
	}

	return client, ""
}

func (s *oauthService) ValidateSession(sessionToken string) (*models.Session, string) {
	session, errCode := s.sessionRepo.GetActiveSessionByToken(sessionToken)
	if errCode != "" {
		return nil, "INVALID_SESSION"
	}

	// Update last activity
	if errCode := s.sessionRepo.UpdateSessionActivity(session.ID); errCode != "" {
		return nil, "SERVER_ERROR"
	}

	return session, ""
}

func (s *oauthService) GenerateAuthCodeWithSession(userID, clientID string) (*models.AuthCode, string) {
	authCode := models.AuthCode{
		Code:      utils.GenerateRandomString(32), // implement this helper function
		UserID:    userID,
		ClientID:  clientID,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}

	if errCode := s.authCodeRepo.CreateAuthCode(&authCode); errCode != "" {
		return nil, "SERVER_ERROR"
	}

	return &authCode, ""
}

func (s *oauthService) ValidateAuthCode(code, clientID string) (*models.AuthCode, string) {
	authCode, errCode := s.authCodeRepo.GetAuthCodeByCode(code)
	if errCode != "" {
		return nil, "INVALID_CODE"
	}

	if authCode.ClientID != clientID {
		return nil, "INVALID_CLIENT"
	}

	if time.Now().After(authCode.ExpiresAt) {
		return nil, "CODE_EXPIRED"
	}

	return authCode, ""
}

func (s *oauthService) ExchangeAuthCodeForToken(code, clientID, clientSecret string) (map[string]interface{}, string) {
	// Validate client
	client, errCode := s.clientRepo.GetClientByID(clientID)
	if errCode != "" || client.Secret != clientSecret {
		return nil, "INVALID_CLIENT"
	}

	// Validate auth code
	authCode, errCode := s.ValidateAuthCode(code, clientID)
	if errCode != "" {
		return nil, errCode
	}

	// Generate token
	token, err := s.jwtService.GenerateToken(authCode.UserID)
	if err != "" {
		return nil, "TOKEN_ERROR"
	}

	// Delete used auth code
	if errCode := s.authCodeRepo.DeleteAuthCode(code); errCode != "" {
		return nil, "SERVER_ERROR"
	}

	return map[string]interface{}{
		"access_token": token,
		"token_type":   "Bearer",
		"expires_in":   3600000,
	}, ""
}

func (s *oauthService) ValidateAccessToken(tokenString string) (map[string]string, string) {
	claims, err := s.jwtService.ValidateToken(tokenString)
	if err != "" {
		return nil, "INVALID_TOKEN"
	}

	userID := claims["user_id"].(string)
	return map[string]string{
		"user_id": userID,
		"valid":   "true",
	}, ""
}
