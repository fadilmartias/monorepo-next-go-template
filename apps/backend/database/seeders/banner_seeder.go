package seeders

import (
	"fmt"
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"

	"gorm.io/gorm"
)

func SeedBanner(db *gorm.DB, count int) {
	log.Printf("Seeding %d banner...", count)
	var items []models.Banner

	// Optional: predefined example
	sample := models.Banner{
		// TODO: fill default values here
	}
	items = append(items, sample)

	// for i := 0; i < count; i++ {
	// 	item := factories.NewBanner()
	// 	// Optional: modify item before hashing/storing
	// 	items = append(items, item)
	// }

	result := db.CreateInBatches(&items, 100)
	if result.Error != nil {
		log.Printf("Could not seed banner: %v", result.Error)
	} else {
		fmt.Printf("Seeded %d banner successfully.\n", count)
	}
}
