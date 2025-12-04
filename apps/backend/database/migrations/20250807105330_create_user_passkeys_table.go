package migrations

import (
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/app/models"
	"gorm.io/gorm"
)

func init() {
	RegisterMigration("create_user_passkeys_table", Up_20250807105330_create_user_passkeys_table, Down_20250807105330_create_user_passkeys_table)
}

func Up_20250807105330_create_user_passkeys_table(db *gorm.DB) {
	log.Println("Running migration: create_user_passkeys_table (UP)")

	err := db.AutoMigrate(&models.UserPasskey{})
	if err != nil {
		log.Fatalf("Could not create user passkeys table: %v", err)
	}

	log.Println("Migration create_user_passkeys_table completed successfully.")
}

func Down_20250807105330_create_user_passkeys_table(db *gorm.DB) {
	log.Println("Running migration: create_user_passkeys_table (DOWN)")

	if db.Migrator().HasTable(&models.UserPasskey{}) {
		err := db.Migrator().DropTable(&models.UserPasskey{})
		if err != nil {
			log.Fatalf("Could not drop user passkeys table: %v", err)
		}
	}
	log.Println("Rollback create_user_passkeys_table completed successfully.")
}
