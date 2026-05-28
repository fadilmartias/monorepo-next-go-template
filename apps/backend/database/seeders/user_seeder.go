package seeders

import (
	"fmt"
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB, count int) {
	log.Printf("Seeding %d users...", count)
	users := make([]models.User, 0, count)
	admin := models.User{
		BaseModelWithDeletedAt: models.BaseModelWithDeletedAt{ID: "A1"},
		Name:                   "Admin",
		Phone:                  "08123456789",
		Email:                  "admin@gmail.com",
		Password:               "namakau123",
		Role:                   "superadmin",
		CurrentBalance:         10000,
	}
	if err := admin.HashPassword(admin.Password); err != nil {
		log.Printf("Failed to hash password for admin: %v", err)
	}
	users = append(users, admin)
	user := models.User{
		BaseModelWithDeletedAt: models.BaseModelWithDeletedAt{ID: "A2"},
		Name:                   "User",
		Phone:                  "08123456780",
		Email:                  "user@gmail.com",
		Password:               "namakau123",
		Role:                   "user",
		CurrentBalance:         0,
	}
	if err := user.HashPassword(user.Password); err != nil {
		log.Printf("Failed to hash password for user: %v", err)
	}
	users = append(users, user)
	result := db.CreateInBatches(&users, 100)
	if result.Error != nil {
		log.Printf("Could not seed user: %v", result.Error)
	}
	fmt.Printf("Seeded %d users successfully.\n", count)
}
