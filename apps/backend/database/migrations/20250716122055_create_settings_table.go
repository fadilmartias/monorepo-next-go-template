package migrations

import (
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration("create_settings_table", Up_20250716122055_create_settings_table, Down_20250716122055_create_settings_table)
}

func Up_20250716122055_create_settings_table(db *gorm.DB) {
	log.Println("Running migration: create_settings_table (UP)")

	if !db.Migrator().HasTable(&models.Setting{}) {
		err := db.Migrator().CreateTable(&models.Setting{})
		if err != nil {
			log.Fatalf("Could not create setting table: %v", err)
		}
	}

	log.Println("Migration create_settings_table completed successfully.")
}

func Down_20250716122055_create_settings_table(db *gorm.DB) {
	log.Println("Running migration: create_settings_table (DOWN)")

	if db.Migrator().HasTable(&models.Setting{}) {
		err := db.Migrator().DropTable(&models.Setting{})
		if err != nil {
			log.Fatalf("Could not drop setting table: %v", err)
		}
	}
	log.Println("Rollback create_settings_table completed successfully.")
}
