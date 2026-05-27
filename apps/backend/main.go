package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/jobs"
	job_handlers "github.com/fadilmartias/dilz_code/apps/backend/app/jobs/handlers"
	"github.com/fadilmartias/dilz_code/apps/backend/app/logger" // Import logger kamu
	"github.com/fadilmartias/dilz_code/apps/backend/bootstrap"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
)

func main() {
	app, db, redis := bootstrap.NewApp()
	jobs.InitQueue(redis)

	// 1. Jalankan Asynq Server (Start() sudah non-blocking dari asalnya)
	logger.Debug("Starting Asynq server...")
	if err := jobs.AsynqServer.Start(job_handlers.NewHandler(db, redis)); err != nil {
		logger.Errorf("Asynq server error: %v", err)
	}
	logger.Debug("Asynq server is ready")

	appConfig := config.LoadAppConfig()

	// 2. Buat channel untuk mendengarkan sinyal OS (Ctrl+C, Docker stop, dll)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// 3. Jalankan Fiber di dalam goroutine agar main thread tidak terblokir
	go func() {
		logger.Debugf("Starting Fiber server on port %s", appConfig.Port)
		if err := app.Listen(appConfig.Port); err != nil {
			logger.Errorf("Fiber server error: %v", err)
		}
	}()

	// 4. Thread utama akan berhenti di sini dan menunggu sinyal OS masuk
	<-quit
	logger.Debug("Menerima sinyal mati (SIGINT/SIGTERM). Memulai Graceful Shutdown...")

	// ==========================================
	// SEQUENCE SHUTDOWN YANG BENAR (BEST PRACTICE)
	// ==========================================

	// Tahap 1: Matikan Fiber dengan timeout (Tolak request HTTP baru, selesaikan yg lagi jalan)
	logger.Debug("Menutup server Fiber...")
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		logger.Errorf("Fiber shutdown error: %v", err)
	}

	// Tahap 2: Matikan Asynq Server (Stop processing background jobs)
	logger.Debug("Menutup server Asynq...")
	jobs.AsynqServer.Stop()

	// Tahap 3: Tutup koneksi Database (opsional tapi dianjurkan)
	if sqlDB, err := db.DB(); err == nil {
		logger.Debug("Menutup koneksi Database...")
		sqlDB.Close()
	}

	// Tahap 4: Sync Logger terakhir (Wajib!)
	logger.Debug("Menyimpan sisa log ke file...")
	logger.Sync()

	logger.Debug("Aplikasi berhasil ditutup dengan aman. Goodbye!")
}
