package main

import (
	"app-one-backend/handlers"
	"app-one-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// CORS setup
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Public routes
	r.POST("/auth/callback", handlers.HandleAuthCallback)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("/me", handlers.GetMe)
	}

	r.Run(":8081")

}
