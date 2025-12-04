package migrations

import (
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration("add_column_total_exp_user_table", Up_20250716130300_add_column_total_exp_user_table, Down_20250716130300_add_column_total_exp_user_table)
}

func Up_20250716130300_add_column_total_exp_user_table(db *gorm.DB) {
	log.Println("Running migration: add_column_total_exp_user_table (UP)")

	err := db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Could not create total_exp column: %v", err)
	}

	log.Println("Migration add_column_total_exp_user_table completed successfully.")
}

func Down_20250716130300_add_column_total_exp_user_table(db *gorm.DB) {
	log.Println("Running migration: add_column_total_exp_user_table (DOWN)")

	if db.Migrator().HasColumn(&models.User{}, "total_exp") {
		err := db.Migrator().DropColumn(&models.User{}, "total_exp")
		if err != nil {
			log.Fatalf("Could not drop total_exp column: %v", err)
		}
	}
	log.Println("Rollback add_column_total_exp_user_table completed successfully.")
}
