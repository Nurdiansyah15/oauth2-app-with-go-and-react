package seeders

import (
	"auth-server/database"
	"auth-server/models"
	"log"
)

func UserSeeder() {
	// TODO: implement user seeder
	users := []models.User{
		{
			Username: "Nurdiansyah",
			Password: "nurdiansyah",
		},
		{
			Username: "Muhammad Nurdiansyah",
			Password: "nurdiansyah",
		},
	}

	for _, user := range users {
		if err := database.DB.FirstOrCreate(&user, models.User{Username: user.Username}).Error; err != nil {
			log.Printf("Error seeding user: %v", err)
		}
	}
}
