package main

import (
	"auth-server/database"
	"auth-server/middleware"
	"auth-server/routes"
	"auth-server/seeders"
	"fmt"
	"html/template"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables from the system.")
	}

	// Initialize DB
	if err := database.InitDB(); err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	seeders.Seeder()


	r := gin.Default()

	// CORS setup
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://localhost:8081",
		"http://localhost:8082",
	}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Security middleware
	r.Use(middleware.SecurityMiddleware())

	r.Static("/assets", "./public/assets")

	// Load HTML templates
	r.SetHTMLTemplate(template.Must(template.New("").ParseGlob("./templates/*")))

	store := cookie.NewStore([]byte(os.Getenv("SESSION_KEY")))
	r.Use(sessions.Sessions("auth-session", store))

	routes.InitUserRoutes(database.DB, r)
	routes.InitOauthRoutes(database.DB, r)

	r.Run(":8080")
}
