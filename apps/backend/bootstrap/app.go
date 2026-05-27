package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fadilmartias/dilz_code/apps/backend/app/client"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/fadilmartias/dilz_code/apps/backend/app/http/middleware"
	"github.com/fadilmartias/dilz_code/apps/backend/app/logger"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	"github.com/fadilmartias/dilz_code/apps/backend/config"
	"github.com/fadilmartias/dilz_code/apps/backend/cronjob"
	"github.com/fadilmartias/dilz_code/apps/backend/routes"

	"github.com/bytedance/sonic"

	"github.com/coregx/coregex"
	"github.com/go-redis/redis/v8"
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
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func NewApp() (*fiber.App, *gorm.DB, *redis.Client) {

	// Init logger
	logger.Init()

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		logger.Debug("Could not load .env file")
	}

	// Create app
	appConfig := config.LoadAppConfig()
	redisStorageNamespace := appConfig.Name
	if redisStorageNamespace == "" {
		redisStorageNamespace = "backend"
	}
	redisStorage := client.NewFiberRedisStorage(nil, redisStorageNamespace)
	redisClient, err := client.ConnectRedis()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	redisStorage.SetClient(redisClient)

	app := fiber.New(fiber.Config{
		AppName:           appConfig.Name,
		SharedStorage:     redisStorage,
		SharedStatePrefix: appConfig.Name + "-shared-",
		JSONEncoder:       sonic.Marshal,
		JSONDecoder:       sonic.Unmarshal,
		RegexHandler:      coregex.MustCompile,
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
		Format:     "[${time}] ${status} | ${latency} | ${method} ${path} | req_id: ${locals:requestid} | error: ${error}\n",
		TimeFormat: "02/01/2006 15:04:05",
	}))

	// DB connection
	db, err := client.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	redis := redisClient

	// Subscribe di startup
	pubSub := redis.Subscribe(context.Background(), "leaderboard-global")
	go func() {
		defer pubSub.Close()
		for msg := range pubSub.Channel() {
			fmt.Println("Dapet update leaderboard:", msg.Payload)

			// Broadcast ke websocket clients
			utils.WebsocketBroadcast("leaderboard-global", msg.Payload)
		}
	}()

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
	app.Use("/statics", static.New("./public"))
	app.Use("/uploads", static.New("./storage/uploads"))
	app.Get("/metrics", monitor.New(monitor.Config{Title: "Firavel Metrics Page"}))
	app.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	cronjob.StartCronJob(db)

	// Register routes
	routes.RegisterApiRoutes(app, db, redis)
	routes.RegisterWebsocketRoutes(app)

	return app, db, redis
}
