package migrations

import (
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration("create_payment_methods_table", Up_20250707194918_create_payment_methods_table, Down_20250707194918_create_payment_methods_table)
}

func Up_20250707194918_create_payment_methods_table(db *gorm.DB) {
	log.Println("Running migration: create_payment_methods_table (UP)")

	err := db.AutoMigrate(&models.PaymentMethod{})
	if err != nil {
		log.Fatalf("Could not create payment method table: %v", err)
	}

	log.Println("Migration create_payment_methods_table completed successfully.")
}

func Down_20250707194918_create_payment_methods_table(db *gorm.DB) {
	log.Println("Running migration: create_payment_methods_table (DOWN)")

	if db.Migrator().HasTable(&models.PaymentMethod{}) {
		err := db.Migrator().DropTable(&models.PaymentMethod{})
		if err != nil {
			log.Fatalf("Could not drop payment method table: %v", err)
		}
	}
	log.Println("Rollback create_payment_methods_table completed successfully.")
}
