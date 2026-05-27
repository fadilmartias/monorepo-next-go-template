package main

import (
	"log"
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/jobs"
	job_handlers "github.com/fadilmartias/dilz_code/apps/backend/app/jobs/handlers"
	"github.com/fadilmartias/dilz_code/apps/backend/bootstrap"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
)

func main() {
	// Buat instance aplikasi dari bootstrap
	app, db, redis := bootstrap.NewApp()
	jobs.InitQueue(redis)
	// Channel untuk tunggu Asynq ready
	ready := make(chan struct{})
	go func() {
		log.Println("Starting Asynq server...")
		if err := jobs.AsynqServer.Start(job_handlers.NewHandler(db, redis)); err != nil {
			log.Fatalf("Asynq server error: %v", err)
		}
		close(ready)
	}()

	// Tunggu Asynq ready
	select {
	case <-ready:
		log.Println("Asynq server is ready")
	case <-time.After(10 * time.Second):
		log.Fatal("Asynq server failed to start within 10 seconds")
	}

	// Muat konfigurasi port dari config
	appConfig := config.LoadAppConfig()

	app.Hooks().OnPostShutdown(func(err error) error {
		if err != nil {
			log.Printf("Shutdown error: %v", err)
		} else {
			log.Println("Shutdown successful")
		}
		return nil
	})
	app.Listen(appConfig.Port)

}
