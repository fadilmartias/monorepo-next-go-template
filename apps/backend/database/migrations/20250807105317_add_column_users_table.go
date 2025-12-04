package migrations

import (
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(
		"add_columns_to_users_table",
		Up_20250807105315_add_columns_to_users_table,
		Down_20250807105315_add_columns_to_users_table,
	)
}

func Up_20250807105315_add_columns_to_users_table(db *gorm.DB) {
	log.Println("Running migration: add_columns_to_users_table (UP)")

	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Could not create user table: %v", err)
	}

	log.Println("Migration create_users_table completed successfully.")
}

func Down_20250807105315_add_columns_to_users_table(db *gorm.DB) {
	log.Println("Running migration: add_columns_to_users_table (DOWN)")

	columns := []string{
		"exp",
		"level",
		"total_spent",
		"avatar_id",
		"border_id",
		"background_id",
		"title_id",
	}

	for _, col := range columns {
		if db.Migrator().HasColumn(&models.User{}, col) {
			if err := db.Migrator().DropColumn(&models.User{}, col); err != nil {
				log.Printf("❌ Failed to drop column %s: %v", col, err)
			} else {
				log.Printf("✅ Column %s dropped.", col)
			}
		} else {
			log.Printf("⚠️ Column %s not found, skipping.", col)
		}
	}
}
