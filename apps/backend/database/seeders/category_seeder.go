package seeders

import (
	"fmt"
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

type categorySeed struct {
	ID       string
	ParentID *string
	Order    int
	Name     string
	Desc     *string
	Slug     string
	IsActive bool
	Type     models.CategoryType
}

func SeedCategory(db *gorm.DB, count int) {
	var seeds = []categorySeed{
		{
			ID:       "C1",
			Name:     "News & Promo",
			Slug:     "news-promo",
			Type:     models.CategoryTypeArticle,
			IsActive: true,
		},
	}

	items := make([]models.Category, 0, len(seeds))
	for _, seed := range seeds {
		items = append(items, models.Category{
			BaseModel: models.BaseModel{ID: seed.ID},
			ParentID:  seed.ParentID,
			Order:     seed.Order,
			Name:      seed.Name,
			Desc:      seed.Desc,
			Slug:      seed.Slug,
			IsActive:  seed.IsActive,
			Type:      seed.Type,
		})
	}

	result := db.Create(&items)
	if result.Error != nil {
		log.Printf("Could not seed category: %v", result.Error)
	} else {
		fmt.Printf("Seeded %d category successfully.\n", count)
	}
}
