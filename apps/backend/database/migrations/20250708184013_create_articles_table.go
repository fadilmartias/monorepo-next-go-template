package migrations

import (
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration("create_articles_table", Up_20250708184013_create_articles_table, Down_20250708184013_create_articles_table)
}

func Up_20250708184013_create_articles_table(db *gorm.DB) {
	log.Println("Running migration: create_articles_table (UP)")

	if !db.Migrator().HasTable(&models.Article{}) {
		err := db.Migrator().CreateTable(&models.Article{})
		if err != nil {
			log.Fatalf("Could not create article table: %v", err)
		}
	}

	log.Println("Migration create_articles_table completed successfully.")
}

func Down_20250708184013_create_articles_table(db *gorm.DB) {
	log.Println("Running migration: create_articles_table (DOWN)")

	if db.Migrator().HasTable(&models.Article{}) {
		err := db.Migrator().DropTable(&models.Article{})
		if err != nil {
			log.Fatalf("Could not drop article table: %v", err)
		}
	}
	log.Println("Rollback create_articles_table completed successfully.")
}
