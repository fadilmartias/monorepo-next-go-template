package seeders

import (
	"fmt"
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"

	"gorm.io/gorm"
)

func SeedCategory(db *gorm.DB, count int) {
	mlId := "DF"
	genshinId := "DG"
	hokId := "DH"
	hsrId := "DI"
	var items = []models.Category{
		{
			ID:       "C1",
			Name:     "Pulsa",
			Slug:     "pulsa",
			Type:     "product",
			IsActive: true,
		},
		{
			ID:       "C2",
			Name:     "Paket Data",
			Slug:     "paket-data",
			Type:     "product",
			IsActive: true,
		},
		{
			ID:       "C4",
			Name:     "Voucher",
			Slug:     "voucher",
			Type:     "product",
			IsActive: true,
		},
		{
			ID:       "C5",
			Name:     "Listrik & Air",
			Slug:     "listrik-air",
			Type:     "product",
			IsActive: true,
		},
		{
			ID:       "C6",
			Name:     "BPJS",
			Slug:     "bpjs",
			Type:     "product",
			IsActive: true,
		},
		{
			ID:       "C7",
			Name:     "Telkom",
			Slug:     "telkom",
			Type:     "product",
			IsActive: true,
		},
		{
			ID:       "C8",
			Name:     "Speedy",
			Slug:     "speedy",
			Type:     "product",
			IsActive: true,
		},
		{
			ID:       "C9",
			Name:     "Indihome",
			Slug:     "indihome",
			Type:     "product",
			IsActive: true,
		},
		{
			ID:       "D0",
			Name:     "Internet",
			Slug:     "internet",
			Type:     "product",
			IsActive: true,
		},
		{
			ID:       "D1",
			Name:     "TV Kabel",
			Slug:     "tv-kabel",
			Type:     "product",
			IsActive: true,
		},
		{
			ID:       "D2",
			Name:     "Telpon",
			Slug:     "telpon",
			Type:     "product",
			IsActive: true,
		},
		{
			ID:       "D3",
			Name:     "Gas",
			Slug:     "gas",
			Type:     "product",
			IsActive: true,
		},
		{
			ID:       "D4",
			Name:     "PDAM",
			Slug:     "pdam",
			Type:     "product",
			IsActive: true,
		},
		{
			ID:       "D5",
			Name:     "Streaming",
			Slug:     "streaming",
			Type:     "product",
			IsActive: true,
		},
		{
			ID:       "D6",
			Name:     "Flash Sale",
			Slug:     "flash-sale",
			Type:     "product",
			IsActive: true,
		},
		{
			ID:       "D7",
			Name:     "Game",
			Slug:     "game",
			Type:     "product",
			IsActive: true,
		},
		{
			ID:       "DA",
			Order:    3,
			Name:     "Virtual Account",
			Slug:     "virtual-account",
			Type:     "payment-method",
			IsActive: true,
		},
		{
			ID:       "DB",
			Order:    2,
			Name:     "E-Wallet",
			Slug:     "e-wallet",
			Type:     "payment-method",
			IsActive: true,
		},
		{
			ID:       "DC",
			Order:    6,
			Name:     "Bank Transfer",
			Slug:     "bank-transfer",
			Type:     "payment-method",
			IsActive: true,
		},
		{
			ID:       "DD",
			Name:     "QRIS",
			Order:    1,
			Slug:     "qris",
			Type:     "payment-method",
			IsActive: true,
		},
		{
			ID:       "DE",
			Order:    5,
			Name:     "Kartu",
			Slug:     "kartu",
			Type:     "payment-method",
			IsActive: true,
		},
		{
			ID:       "D10",
			Order:    4,
			Name:     "Convenient Store",
			Slug:     "convenient-store",
			Type:     "payment-method",
			IsActive: true,
		},
		{
			ID:       mlId,
			Name:     "Mobile Legend",
			Slug:     "mobile-legend",
			Type:     "product-parent",
			IsActive: true,
		},
		{
			ID:       genshinId,
			Name:     "Genshin Impact",
			Slug:     "genshin-impact",
			Type:     "product-parent",
			IsActive: true,
		},
		{
			ID:       hokId,
			Name:     "Honor of Kings",
			Slug:     "honor-of-kings",
			Type:     "product-parent",
			IsActive: true,
		},
		{
			ID:       hsrId,
			Name:     "Honkai Star Rail",
			Slug:     "honkai-star-rail",
			Type:     "product-parent",
			IsActive: true,
		},
		{
			ID:       "D8",
			ParentID: &mlId,
			Order:    1,
			Name:     "Weekly Diamond Pass",
			Slug:     "weekly-diamond-pass",
			Type:     "product-subcategory",
			IsActive: true,
		},
		{
			ID:       "D9",
			ParentID: &mlId,
			Order:    2,
			Name:     "Diamond",
			Slug:     "diamond",
			Type:     "product-subcategory",
			IsActive: true,
		},
		{
			ID:       "D11",
			ParentID: &hokId,
			Order:    1,
			Name:     "Weekly Card",
			Slug:     "weekly-card",
			Type:     "product-subcategory",
			IsActive: true,
		},
		{
			ID:       "D12",
			ParentID: &hokId,
			Order:    2,
			Name:     "Token",
			Slug:     "token",
			Type:     "product-subcategory",
			IsActive: true,
		},
		{
			ID:       "D13",
			ParentID: &genshinId,
			Order:    1,
			Name:     "Blessing of the Welkin Moon",
			Slug:     "blessing-of-the-welkin-moon",
			Type:     "product-subcategory",
			IsActive: true,
		},
		{
			ID:       "D14",
			ParentID: &genshinId,
			Order:    2,
			Name:     "Genesis Crystal",
			Slug:     "genesis-crystal",
			Type:     "product-subcategory",
			IsActive: true,
		},
		{
			ID:       "D15",
			ParentID: &hsrId,
			Order:    1,
			Name:     "Express Supply Pass",
			Slug:     "express-supply-pass",
			Type:     "product-subcategory",
			IsActive: true,
		},
		{
			ID:       "D16",
			ParentID: &hsrId,
			Order:    2,
			Name:     "Oneiric Shard",
			Slug:     "oneiric-shard",
			Type:     "product-subcategory",
			IsActive: true,
		},
		{
			ID:       "D17",
			Name:     "Berita Game",
			Slug:     "berita-game",
			Type:     "article",
			IsActive: true,
		},
		{
			ID:       "D18",
			Name:     "Voucher Robux",
			Slug:     "voucher-robux",
			Type:     "product-subcategory",
			IsActive: true,
		},
	}

	result := db.Create(&items)
	if result.Error != nil {
		log.Printf("Could not seed category: %v", result.Error)
	} else {
		fmt.Printf("Seeded %d category successfully.\n", count)
	}
}
