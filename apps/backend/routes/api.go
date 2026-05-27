package routes

import (
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	controllers_v0 "github.com/fadilmartias/dilz_code/apps/backend/app/http/controllers/v0"
	controllers_v1 "github.com/fadilmartias/dilz_code/apps/backend/app/http/controllers/v1"
	"github.com/fadilmartias/dilz_code/apps/backend/app/http/middleware"
	"github.com/fadilmartias/dilz_code/apps/backend/app/repositories"
	"github.com/fadilmartias/dilz_code/apps/backend/app/requests"
	"github.com/fadilmartias/dilz_code/apps/backend/app/services"
	"github.com/fadilmartias/dilz_code/apps/backend/app/usecases"
	"github.com/fadilmartias/dilz_code/apps/backend/app/utils"
	_ "github.com/fadilmartias/dilz_code/apps/backend/docs"
	"github.com/fadilmartias/dilz_code/apps/backend/graph"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"gorm.io/gorm"
)

func RegisterApiRoutes(app *fiber.App, db *gorm.DB, redis *redis.Client) {
	// ========= GLOBAL MIDDLEWARE =========
	app.Use(middleware.TraceIDMiddleware())
	app.Use(middleware.ActivityContextMiddleware())
	app.Use(middleware.RobotTag())
	app.Use(middleware.GetUser())

	// ========= ROOT ROUTES =========
	app.Get("/", func(c fiber.Ctx) error {
		return utils.SuccessResponse(c, utils.SuccessResponseFormat{
			Message: "Hello, World!",
		})
	})

	app.Get("/swagger.json", func(c fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})

	app.Get("/swagger/*", func(c fiber.Ctx) error {
		html := `<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<title>Swagger UI</title>
			<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.11.0/swagger-ui.min.css" />
			<style>
				body { margin: 0; padding: 0; }
			</style>
		</head>
		<body>
			<div id="swagger-ui"></div>
			<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.11.0/swagger-ui-bundle.js"></script>
			<script>
				window.onload = function() {
					SwaggerUIBundle({
						url: "/swagger.json", // Arahkan ke route JSON yang kita buat di atas
						dom_id: '#swagger-ui',
						presets: [ SwaggerUIBundle.presets.apis ],
					});
				};
			</script>
		</body>
		</html>`

		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(html)
	})

	// ========= API V0 =========
	apiV0 := app.Group("/v0")
	apiV0.Use(middleware.RateLimiter(100, 1*time.Minute))
	apiV0.Use(middleware.Signature())

	genericController := controllers_v0.NewGenericController(db, redis)
	apiV0.Put("/uploads", controllers_v0.NewUploadController(db, redis).Upload).Name("uploads")
	apiV0.Get("/:model", genericController.Index).Name("generic.index")
	apiV0.Get("/:model/:id", genericController.Show).Name("generic.show")
	apiV0.Post("/:model", genericController.Store, middleware.Auth([]string{"admin"}, []string{})).Name("generic.store")
	apiV0.Put("/:model/:id", genericController.Update, middleware.Auth([]string{"admin"}, []string{})).Name("generic.update")
	apiV0.Patch("/:model/:id", genericController.Patch, middleware.Auth([]string{"admin"}, []string{})).Name("generic.patch")
	apiV0.Delete("/:model/:id", genericController.Destroy, middleware.Auth([]string{"admin"}, []string{})).Name("generic.destroy")

	// ========= API V1 =========
	apiV1 := app.Group("/v1")
	apiV1.Use(middleware.RateLimiter(100, 1*time.Minute))
	apiV1.Use(middleware.Signature())
	activityLogRepository := repositories.NewActivityLogRepository(db)

	activityLogService := services.NewActivityLogService(db, activityLogRepository)

	activityLogController := controllers_v1.NewActivityLogController(activityLogService)
	activityLogRoutes := apiV1.Group("/activity-logs")
	{
		activityLogRoutes.Get("/:id", activityLogController.Show).Name("activity-logs.show")
	}

	// --- User Module ---
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(db, redis, userRepository)
	userController := controllers_v1.NewUserController(userService)
	userRoutes := apiV1.Group("/users")
	{
		userRoutes.Get("/", middleware.Auth([]string{"admin"}, []string{}), userController.Index).Name("users.index")
		userRoutes.Get("/:id", middleware.Auth([]string{"admin"}, []string{}), userController.Show).Name("users.show")
		userRoutes.Patch("/profile", middleware.Validator[requests.UpdateProfileInput](), userController.UpdateProfile).Name("users.update-profile")
		userRoutes.Patch("/password", middleware.Validator[requests.UpdatePasswordInput](), userController.UpdatePassword).Name("users.update-password")
	}

	// --- Article Module ---
	articleRepository := repositories.NewArticleRepository(db)
	articleService := services.NewArticleService(db, redis, articleRepository)
	articleController := controllers_v1.NewArticleController(articleService)
	articleRoute := apiV1.Group("/articles")
	{
		articleRoute.Get("/", articleController.Index).Name("articles.index")
		articleRoute.Put("/", articleController.Process, middleware.Auth([]string{"admin"}, []string{})).Name("articles.process")
		articleRoute.Get("/popular", articleController.Popular).Name("articles.popular")
		articleRoute.Get("/:slug", articleController.Show).Name("articles.show")
	}

	// --- Auth Module ---

	emailVerificationUsecase := usecases.NewSendEmailVerificationUsecase()
	emailResetPasswordUsecase := usecases.NewSendResetPasswordEmailUsecase()
	emailService := services.NewEmailService(emailVerificationUsecase, emailResetPasswordUsecase)
	authService := services.NewAuthService(db, redis, userRepository, activityLogService)
	authController := controllers_v1.NewAuthController(db, redis, userService, emailService, authService)
	authRoutes := apiV1.Group("/auth")
	{
		authRoutes.Get("/me", authController.Me).Name("auth.me")
		authRoutes.Post("/logout", authController.Logout).Name("auth.logout")
		authRoutes.Post("/refresh-access-token", authController.RefreshAccessToken).Name("auth.refresh-access-token")
		authRoutes.Post("/login", middleware.Guest(), middleware.Validator[requests.LoginInput](), authController.Login).Name("auth.login")
		authRoutes.Post("/login-2fa", middleware.Guest(), middleware.Validator[requests.Login2FARequest](), authController.Login2FA).Name("auth.login-2fa")
		authRoutes.Get("/google/redirect", middleware.Guest(), authController.GoogleRedirect).Name("auth.google.redirect")
		authRoutes.Get("/google/callback", middleware.Guest(), authController.GoogleCallback).Name("auth.google.callback")
		authRoutes.Post("/google/one-tap", middleware.Guest(), authController.GoogleOneTap).Name("auth.google.one-tap")
		authRoutes.Get("/discord/redirect", middleware.Guest(), authController.DiscordRedirect).Name("auth.discord.redirect")
		authRoutes.Get("/discord/callback", middleware.Guest(), authController.DiscordCallback).Name("auth.discord.callback")
		authRoutes.Get("/facebook/redirect", middleware.Guest(), authController.FacebookRedirect).Name("auth.facebook.redirect")
		authRoutes.Get("/facebook/callback", middleware.Guest(), authController.FacebookCallback).Name("auth.facebook.callback")
		authRoutes.Get("/steam/redirect", middleware.Guest(), authController.SteamRedirect).Name("auth.steam.redirect")
		authRoutes.Get("/steam/callback", middleware.Guest(), authController.SteamCallback).Name("auth.steam.callback")
		authRoutes.Get("/twitch/redirect", middleware.Guest(), authController.TwitchRedirect).Name("auth.twitch.redirect")
		authRoutes.Get("/twitch/callback", middleware.Guest(), authController.TwitchCallback).Name("auth.twitch.callback")
		authRoutes.Post("/register", middleware.RateLimiter(5, 1*time.Minute), middleware.Guest(), middleware.Validator[requests.RegisterInput](), authController.Register).Name("auth.register")
		authRoutes.Post("/forgot-password", middleware.Guest(), middleware.Validator[requests.ForgotPasswordInput](), authController.ForgotPassword).Name("auth.forgot-password")
		authRoutes.Post("/reset-password", middleware.Guest(), middleware.Validator[requests.ResetPasswordInput](), authController.ResetPassword).Name("auth.reset-password")
		authRoutes.Post("/send-email-verification", middleware.Auth([]string{}, []string{}), authController.SendEmailVerification).Name("auth.send-email-verification")
		authRoutes.Post("/verify-email", middleware.Auth([]string{}, []string{}), middleware.Validator[requests.VerifyEmailInput](), authController.VerifyEmail).Name("auth.verify-email")
		authRoutes.Post("/register-2fa", middleware.Auth([]string{}, []string{}), authController.Register2FA).Name("auth.register-2fa")
		authRoutes.Post("/verify-2fa", middleware.Auth([]string{}, []string{}), middleware.Validator[requests.Verify2FARequest](), authController.Verify2FA).Name("auth.verify-2fa")
		authRoutes.Post("/disable-2fa", middleware.Auth([]string{}, []string{}), authController.Disable2FA).Name("auth.disable-2fa")
		authRoutes.Post("/send-otp", middleware.Validator[requests.SendOTPRequest](), authController.SendOTP).Name("auth.send-otp")
		authRoutes.Post("/verify-otp", middleware.Validator[requests.VerifyOTPRequest](), authController.VerifyOTP).Name("auth.verify-otp")
		authRoutes.Post("/verify-recaptcha", authController.VerifyRecaptcha).Name("auth.verify-recaptcha")
	}

	// --- Banner Module ---
	bannerRepository := repositories.NewBannerRepository(db)
	bannerService := services.NewBannerService(db, redis, bannerRepository)
	bannerController := controllers_v1.NewBannerController(bannerService)
	bannerRoute := apiV1.Group("/banners")
	{
		bannerRoute.Get("/", bannerController.Index).Name("banners.index")
		bannerRoute.Put("/", bannerController.Process, middleware.Auth([]string{"admin"}, []string{})).Name("banners.process")
	}

	// --- Product Module ---
	// settingRepository := repositories.NewSettingRepository(db)

	// --- Payment Method Module ---
	paymentMethodRepository := repositories.NewPaymentMethodRepository(db)
	paymentMethodService := services.NewPaymentMethodService(db, redis, paymentMethodRepository)
	paymentMethodController := controllers_v1.NewPaymentMethodController(paymentMethodService)
	paymentMethodRoutes := apiV1.Group("/payment-methods")
	{
		paymentMethodRoutes.Get("/", paymentMethodController.Index).Name("payment-methods.index")
	}

	// --- Setting Module ---
	settingRoutes := apiV1.Group("/settings")
	{
		settingRoutes.Get("/", controllers_v1.NewSettingController(db, redis).Index).Name("settings.index")
		settingRoutes.Get("/:key", controllers_v1.NewSettingController(db, redis).Show).Name("settings.show")
		settingRoutes.Post("/", controllers_v1.NewSettingController(db, redis).Store).Name("settings.store")
		settingRoutes.Put("/:id", controllers_v1.NewSettingController(db, redis).Update).Name("settings.update")
		settingRoutes.Delete("/:id", controllers_v1.NewSettingController(db, redis).Destroy).Name("settings.destroy")
	}

	// --- Midtrans Module ---
	// midtransRoutes := apiV1.Group("/midtrans")
	// midtransController := controllers_v1.NewMidtransController(db, redis, midtransService)
	// {
	// 	midtransRoutes.Get("/qr-gopay/:id", midtransController.RenderQrGopay).Name("midtrans.render-qr-gopay")
	// 	midtransRoutes.Get("/balance", midtransController.Balance).Name("midtrans.balance")
	// 	midtransRoutes.Post("/refund/:reference_id", midtransController.RefundTransaction).Name("midtrans.refund-transaction")
	// 	midtransRoutes.Post("/webhook", midtransController.Webhook).Name("midtrans.webhook")
	// }

	// --- Dashboard Module ---
	dashboardRoutes := apiV1.Group("/dashboard")
	{
		dashboardRoutes.Get("/", middleware.Auth([]string{"admin"}, []string{}), controllers_v1.NewDashboardController(db, redis).Index).Name("dashboard.index")
	}

	// GraphQL server
	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{Resolvers: &graph.Resolver{
				UserService: userService,
			}},
		),
	)
	// GraphQL endpoint
	app.All("/graphql", func(c fiber.Ctx) error {
		handler := fasthttpadaptor.NewFastHTTPHandlerFunc(srv.ServeHTTP)
		handler(c.RequestCtx())
		return nil
	})

	app.Get("/playground", func(c fiber.Ctx) error {
		handler := fasthttpadaptor.NewFastHTTPHandlerFunc(
			playground.Handler("GraphQL Playground", "/graphql"),
		)
		handler(c.RequestCtx()) // ⛔ ini wajib dipanggil!
		return nil
	})
}

// fiber:context-methods migrated
