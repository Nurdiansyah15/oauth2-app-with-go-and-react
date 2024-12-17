package seeders

import (
	"auth-server/database"
	"log"
)

func Seeder() {

	log.Println("Clearing existing data...")

	deleteTable("users")
	deleteTable("clients")
	deleteTable("auth_codes")
	deleteTable("sessions")

	log.Println("Seeding data...")

	ClientSeeder()
	UserSeeder()
}

func deleteTable(tableName string) {
	if err := database.DB.Exec("DELETE FROM " + tableName).Error; err != nil {
		log.Fatalf("Failed to clear %s table: %v", tableName, err)
	}
}
