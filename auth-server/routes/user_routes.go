package routes

import (
	"auth-server/handlers"
	"auth-server/middleware"
	"auth-server/repositories"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitUserRoutes(db *gorm.DB, r *gin.Engine) {

	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		userRepo := repositories.NewUserRepository(db)
		userHandler := handlers.NewUserHandler(userRepo)

		protected.GET("/me", userHandler.GetUserInfoHandler)
	}

}
