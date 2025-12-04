package migrations

import (
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration("create_categories_table", Up_20250707195012_create_categories_table, Down_20250707195012_create_categories_table)
}

func Up_20250707195012_create_categories_table(db *gorm.DB) {
	log.Println("Running migration: create_categories_table (UP)")

	err := db.AutoMigrate(&models.Category{})
	if err != nil {
		log.Fatalf("Could not create category table: %v", err)
	}
	log.Println("Migration create_categories_table completed successfully.")
}

func Down_20250707195012_create_categories_table(db *gorm.DB) {
	log.Println("Running migration: create_categories_table (DOWN)")

	if db.Migrator().HasTable(&models.Category{}) {
		err := db.Migrator().DropTable(&models.Category{})
		if err != nil {
			log.Fatalf("Could not drop category table: %v", err)
		}
	}
	log.Println("Rollback create_categories_table completed successfully.")
}
