package migrations

import (
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration("create_banners_table", Up_20250708183954_create_banners_table, Down_20250708183954_create_banners_table)
}

func Up_20250708183954_create_banners_table(db *gorm.DB) {
	log.Println("Running migration: create_banners_table (UP)")

	if !db.Migrator().HasTable(&models.Banner{}) {
		err := db.Migrator().CreateTable(&models.Banner{})
		if err != nil {
			log.Fatalf("Could not create banner table: %v", err)
		}
	}

	log.Println("Migration create_banners_table completed successfully.")
}

func Down_20250708183954_create_banners_table(db *gorm.DB) {
	log.Println("Running migration: create_banners_table (DOWN)")

	if db.Migrator().HasTable(&models.Banner{}) {
		err := db.Migrator().DropTable(&models.Banner{})
		if err != nil {
			log.Fatalf("Could not drop banner table: %v", err)
		}
	}
	log.Println("Rollback create_banners_table completed successfully.")
}
