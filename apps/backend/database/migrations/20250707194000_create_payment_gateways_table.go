package migrations

import (
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration("create_payment_gateways_table", Up_20250707194934_create_payment_gateways_table, Down_20250707194934_create_payment_gateways_table)
}

func Up_20250707194934_create_payment_gateways_table(db *gorm.DB) {
	log.Println("Running migration: create_payment_gateways_table (UP)")

	err := db.AutoMigrate(&models.PaymentGateway{})
	if err != nil {
		log.Fatalf("Could not create payment gateway table: %v", err)
	}

	log.Println("Migration create_payment_gateways_table completed successfully.")
}

func Down_20250707194934_create_payment_gateways_table(db *gorm.DB) {
	log.Println("Running migration: create_payment_gateways_table (DOWN)")

	if db.Migrator().HasTable(&models.PaymentGateway{}) {
		err := db.Migrator().DropTable(&models.PaymentGateway{})
		if err != nil {
			log.Fatalf("Could not drop payment gateway table: %v", err)
		}
	}
	log.Println("Rollback create_payment_gateways_table completed successfully.")
}
