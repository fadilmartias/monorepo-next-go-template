package seeders

import (
	"fmt"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/fadilmartias/dilz_code/apps/backend/app/models"

	"gorm.io/gorm"
)

func SeedPaymentGateway(db *gorm.DB, count int) {
	log.Printf("Seeding %d payment gateway...", count)
	var items []models.PaymentGateway
	meta, _ := sonic.Marshal(map[string]string{
		"charge_endpoint":  "/v2/charge",
		"capture_endpoint": "/v2/capture",
		"invoice_endpoint": "/v1/invoices",
	})
	metaStr := string(meta)
	// Optional: predefined example
	sample := models.PaymentGateway{
		ID:        "PG1",
		Order:     1,
		Name:      "Midtrans",
		Metadata:  &metaStr,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	items = append(items, sample)

	sample2 := models.PaymentGateway{
		ID:        "PG2",
		Order:     2,
		Name:      "iPaymu",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	items = append(items, sample2)

	// for i := 0; i < count; i++ {
	// 	item := factories.NewPaymentGateway()
	// 	// Optional: modify item before hashing/storing
	// 	items = append(items, item)
	// }

	result := db.CreateInBatches(&items, 100)
	if result.Error != nil {
		log.Printf("Could not seed payment gateway: %v", result.Error)
	} else {
		fmt.Printf("Seeded %d payment gateway successfully.\n", count)
	}
}
