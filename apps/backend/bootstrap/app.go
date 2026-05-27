package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/fadilmartias/dilz_code/apps/backend/app/http/middleware"
	"github.com/fadilmartias/dilz_code/apps/backend/app/logger"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"github.com/fadilmartias/dilz_code/apps/backend/cronjob"
	"github.com/fadilmartias/dilz_code/apps/backend/routes"

	"github.com/bytedance/sonic"

	"github.com/gofiber/contrib/v3/monitor"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	fLogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	glogger "gorm.io/gorm/logger"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewApp() (*fiber.App, *gorm.DB, *config.RedisClient) {

	// Init logger
	logger.Init()

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		logger.Error("Could not load .env file")
	}

	// Create app
	app := fiber.New(fiber.Config{
		AppName:     config.LoadAppConfig().Name,
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's a *fiber.Error
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			message := err.Error()
			if message == "" {
				message = "Internal Server Error"
			}

			return utils.ErrorResponse(ctx, utils.ErrorResponseFormat{
				Code:    code,
				Message: message,
				Details: err,
			})
		},
	})

	// Logger middleware
	app.Use(middleware.Logger())

	app.Use(fLogger.New(fLogger.Config{
		Stream: os.Stdout, // Tulis log ke stdout

		// Format teks biasa (mudah dibaca) yang selaras dengan Zap
		Format:     "[${time}] ${status} | ${latency} | ${method} ${path} | req_id:${locals:requestid} | error:${error}\n",
		TimeFormat: "02/01/2006 15:04:05",
	}))

	// DB connection
	db := ConnectDB()

	// Redis connection
	redis := config.NewRedisClient()

	// Subscribe di startup
	redis.Subscribe(context.Background(), "leaderboard-global", func(payload string) {
		fmt.Println("Dapet update leaderboard:", payload)

		// Broadcast ke websocket clients
		utils.WebsocketBroadcast("leaderboard-global", payload)
	})

	// Use middleware
	app.Use(recover.New(recover.Config{
		EnableStackTrace: config.LoadAppConfig().Env != "production",
	}))
	// app.Use(etag.New())
	app.Use(helmet.New(helmet.Config{
		CrossOriginResourcePolicy: "cross-origin",
	}))
	app.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"*"},
		AllowOrigins:     []string{os.Getenv("FE_URL")},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Forwarded-For", "X-Signature", "X-Timestamp", "X-Tenant-Id", "X-Dev-Key", "X-Idempotency-Key"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Set-Cookie"},
	}))
	app.Use(compress.New(compress.Config{

		Level: compress.LevelBestSpeed, // 1
	}))
	app.Use(pprof.New(pprof.Config{
		Next: func(c fiber.Ctx) bool {
			return config.LoadAppConfig().Env != "production"
		},
	}))
	app.Use(csrf.New(csrf.Config{
		IdleTimeout: 10 * time.Minute,
	}))
	app.Use(requestid.New())
	app.Get("/*", static.New("./public")) // Static file
	app.Get("/metrics", monitor.New(monitor.Config{Title: "Firavel Metrics Page"}))
	app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	cronjob.StartCronJob(db, redis)

	// Register routes
	routes.RegisterApiRoutes(app, db, redis)
	routes.RegisterWebsocketRoutes(app)

	return app, db, redis
}

func ConnectDB() *gorm.DB {
	dbConfig := config.LoadDBConfig()
	appConfig := config.LoadAppConfig()

	// Format DSN untuk MySQL
	// format: "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)

	var gormLogger glogger.Interface
	if appConfig.Env != "production" {
		gormLogger = glogger.Default.LogMode(glogger.Info)
	} else {
		gormLogger = glogger.Default.LogMode(glogger.Error)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Could not get database instance: %v", err)
	}
	if appConfig.Env != "production" {
		sqlDB.SetMaxIdleConns(5)  // cukup 5 idle
		sqlDB.SetMaxOpenConns(10) // max 10 koneksi aktif
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
	} else {
		sqlDB.SetMaxIdleConns(20)           // simpan 20 koneksi siap pakai
		sqlDB.SetMaxOpenConns(200)          // max 200 koneksi aktif
		sqlDB.SetConnMaxLifetime(time.Hour) // recycle tiap 1 jam

	}
	return db
}
