package seeders

import (
	"fmt"
	"log"
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"

	"gorm.io/gorm"
)

func SeedArticle(db *gorm.DB, count int) {
	log.Printf("Seeding %d article...", count)
	var items []models.Article

	// Optional: predefined example
	sample := models.Article{
		Title:       "Judul Artikel",
		TenantID:    "T1",
		CategoryID:  "D17",
		Slug:        "judul-artikel",
		Content:     "Isi artikel",
		Img:         "https://example.com/image.jpg",
		PublishedAt: time.Now(),
		IsActive:    true,
	}
	items = append(items, sample)

	// for i := 0; i < count; i++ {
	// 	item := factories.NewArticle()
	// 	// Optional: modify item before hashing/storing
	// 	items = append(items, item)
	// }

	result := db.CreateInBatches(&items, 100)
	if result.Error != nil {
		log.Printf("Could not seed article: %v", result.Error)
	} else {
		fmt.Printf("Seeded %d article successfully.\n", count)
	}
}
