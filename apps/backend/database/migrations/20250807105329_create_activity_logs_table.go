package migrations

import (
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration("create_activity_logs_table", Up_20250807105329_create_activity_logs_table, Down_20250807105329_create_activity_logs_table)
}

func Up_20250807105329_create_activity_logs_table(db *gorm.DB) {
	log.Println("Running migration: create_activity_logs_table (UP)")

	err := db.AutoMigrate(&models.ActivityLog{})
	if err != nil {
		log.Fatalf("Could not create activity logs table: %v", err)
	}

	log.Println("Migration create_activity_logs_table completed successfully.")
}

func Down_20250807105329_create_activity_logs_table(db *gorm.DB) {
	log.Println("Running migration: create_activity_logs_table (DOWN)")

	if db.Migrator().HasTable(&models.ActivityLog{}) {
		err := db.Migrator().DropTable(&models.ActivityLog{})
		if err != nil {
			log.Fatalf("Could not drop activity logs table: %v", err)
		}
	}
	log.Println("Rollback add_profit_column_product_orders_table completed successfully.")
}
