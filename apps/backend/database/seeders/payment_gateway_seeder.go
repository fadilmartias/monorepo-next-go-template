package seeders

import (
	"fmt"
	"log"

	"github.com/bytedance/sonic"
	"github.com/fadilmartias/dilz_code/apps/backend/app/models"

	"gorm.io/gorm"
)

func newPaymentGatewaySeed(id string, order int, name string, metadata *string, isActive bool) models.PaymentGateway {
	sample := models.PaymentGateway{
		Order:    order,
		Name:     name,
		Metadata: metadata,
		IsActive: isActive,
	}

	sample.ID = id

	return sample
}

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
	sample := newPaymentGatewaySeed("PG1", 1, "Midtrans", &metaStr, true)
	items = append(items, sample)

	sample2 := newPaymentGatewaySeed("PG2", 2, "iPaymu", nil, true)
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
