package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
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

	// Tambahkan handler shutdown di goroutine
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Shutting down...")
		jobs.AsynqClient.Close()
		jobs.AsynqServer.Stop()
		redis.Close()
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Shutdown error: %v", err)
		}
	}()

	// Jalankan server
	log.Fatal(app.Listen(appConfig.Port))
}
