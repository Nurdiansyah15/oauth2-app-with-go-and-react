package routes

import (
	"auth-server/handlers"
	"auth-server/repositories"
	"auth-server/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitOauthRoutes(db *gorm.DB, r *gin.Engine) {

	clientRepo := repositories.NewClientRepository(db)
	sessionRepo := repositories.NewSessionRepository(db)
	authCodeRepo := repositories.NewAuthCodeRepository(db)
	userRepo := repositories.NewUserRepository(db)
	jwtService := services.NewJWTService()
	authService := services.NewAuthService(userRepo)
	oauthService := services.NewOAuthService(clientRepo, sessionRepo, authCodeRepo, jwtService)
	oauthHandler := handlers.NewOauthHandler(clientRepo, sessionRepo, authCodeRepo, authService, oauthService, jwtService)

	r.GET("/login", oauthHandler.LoginPage)
	r.POST("/oauth/login", oauthHandler.LoginHandler)
	r.GET("/oauth/authorize", oauthHandler.AuthorizeHandler)
	r.POST("/oauth/token", oauthHandler.TokenHandler)
	r.GET("/oauth/validate", oauthHandler.ValidateTokenHandler)
	r.GET("/oauth/logout", oauthHandler.LogoutHandler)
}
