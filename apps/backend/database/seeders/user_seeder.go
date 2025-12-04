package seeders

import (
	"fmt"
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	_ "github.com/fadilmartias/dilz_code/apps/backend/app/models"

	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB, count int) {
	log.Printf("Seeding %d users...", count)
	users := make([]models.User, 0, count)
	admin := models.User{
		ID:       "A1",
		TenantID: "T1",
		Name:     "Admin",
		Email:    "admin@gmail.com",
		Phone:    "08123456789",
		Password: "namakau123",
		Role:     "admin",
	}
	if err := admin.HashPassword(admin.Password); err != nil {
		log.Printf("Failed to hash password for admin: %v", err)
	}
	users = append(users, admin)
	user := models.User{
		ID:       "A2",
		TenantID: "T1",
		Name:     "User",
		Email:    "user@gmail.com",
		Phone:    "08123456780",
		Password: "namakau123",
		Role:     "user",
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
