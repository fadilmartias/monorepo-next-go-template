package seeders

import (
	"log"

	"github.com/fadilmartias/dilz_code/apps/backend/bootstrap"
)

// RunAllSeeders adalah entry point untuk menjalankan semua seeder.
func RunAllSeeders() {
	db, err := bootstrap.ConnectDB()
	if err != nil {
		panic(err)
	}
	log.Println("Running all seeders...")

	SeedUsers(db, 20)          // Buat 20 user palsu
	SeedCategory(db, 20)       // Buat 20 user palsu
	SeedPaymentGateway(db, 20) // Buat 20 user palsu
	SeedPaymentMethod(db, 20)  // Buat 20 user palsu
	SeedSetting(db, 20)        // Buat 20 user palsu

	log.Println("All seeders completed.")
}
