package seeders

import (
	"auth-server/database"
	"auth-server/models"
	"log"
)

func ClientSeeder() {
	// TODO: implement client seeder
	clients := []models.Client{
		{
			ID:          "app-one-client",
			Secret:      "app-one-secret",
			RedirectURI: "http://localhost:3000/callback",
		},
		{
			ID:          "app-two-client",
			Secret:      "app-two-secret",
			RedirectURI: "http://localhost:3001/callback",
		},
	}

	for _, client := range clients {
		if err := database.DB.FirstOrCreate(&client, models.Client{Secret: client.Secret}).Error; err != nil {
			log.Printf("Error seeding client: %v", err)
		}
	}

	log.Println("Database seeding completed!")
}
