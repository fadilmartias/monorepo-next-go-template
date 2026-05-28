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
			ID:       "DA",
			Order:    3,
			Name:     "Virtual Account",
			Slug:     "virtual-account",
			Type:     models.CategoryTypePaymentMethod,
			IsActive: true,
		},
		{
			ID:       "DB",
			Order:    2,
			Name:     "E-Wallet",
			Slug:     "e-wallet",
			Type:     models.CategoryTypePaymentMethod,
			IsActive: true,
		},
		{
			ID:       "DC",
			Order:    6,
			Name:     "Bank Transfer",
			Slug:     "bank-transfer",
			Type:     models.CategoryTypePaymentMethod,
			IsActive: true,
		},
		{
			ID:       "DD",
			Name:     "QRIS",
			Order:    1,
			Slug:     "qris",
			Type:     models.CategoryTypePaymentMethod,
			IsActive: true,
		},
		{
			ID:       "DE",
			Order:    5,
			Name:     "Kartu",
			Slug:     "kartu",
			Type:     models.CategoryTypePaymentMethod,
			IsActive: true,
		},
		{
			ID:       "D10",
			Order:    4,
			Name:     "Convenient Store",
			Slug:     "convenient-store",
			Type:     models.CategoryTypePaymentMethod,
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
