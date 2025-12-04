package migrations

import (
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration("add_column_totp_secret_user_table", Up_20250807105315_add_column_totp_secret_user_table, Down_20250807105315_add_column_totp_secret_user_table)
}

func Up_20250807105315_add_column_totp_secret_user_table(db *gorm.DB) {
	log.Println("Running migration: add_column_totp_secret_user_table (UP)")

	if !db.Migrator().HasColumn(&models.User{}, "totp_secret") {
		err := db.Migrator().AddColumn(&models.User{}, "totp_secret")
		if err != nil {
			log.Fatalf("Failed to add column totp_secret: %v", err)
		}
		log.Println("Column totp_secret added.")
	} else {
		log.Println("Column totp_secret already exists, skipping.")
	}
}

func Down_20250807105315_add_column_totp_secret_user_table(db *gorm.DB) {
	log.Println("Running migration: add_column_totp_secret_user_table (DOWN)")

	if db.Migrator().HasColumn(&models.User{}, "totp_secret") {
		err := db.Migrator().DropColumn(&models.User{}, "totp_secret")
		if err != nil {
			log.Fatalf("Failed to drop column totp_secret: %v", err)
		}
		log.Println("Column totp_secret dropped.")
	} else {
		log.Println("Column totp_secret not found, skipping.")
	}
}
