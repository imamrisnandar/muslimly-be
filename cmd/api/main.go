package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/time/rate"

	_ "muslimly-be/docs" // Import generated docs
	"muslimly-be/internal/api/router"

	authHandler "muslimly-be/internal/features/auth/handler"
	authService "muslimly-be/internal/features/auth/service"

	userHandler "muslimly-be/internal/features/user/handler"
	userModel "muslimly-be/internal/features/user/model"
	userRepo "muslimly-be/internal/features/user/repository"
	userService "muslimly-be/internal/features/user/service"

	syncHandler "muslimly-be/internal/features/sync/handler"
	syncModel "muslimly-be/internal/features/sync/model"
	syncRepo "muslimly-be/internal/features/sync/repository"
	syncService "muslimly-be/internal/features/sync/service"

	userSettingsHandler "muslimly-be/internal/features/user_settings/handler"
	userSettingsModel "muslimly-be/internal/features/user_settings/model"
	userSettingsRepo "muslimly-be/internal/features/user_settings/repository"
	userSettingsService "muslimly-be/internal/features/user_settings/service"

	notifHandler "muslimly-be/internal/features/notification/handler"
	notifModel "muslimly-be/internal/features/notification/model"
	notifRepo "muslimly-be/internal/features/notification/repository"
	notifService "muslimly-be/internal/features/notification/service"

	appConfigHandler "muslimly-be/internal/features/app_config/handler"
	appConfigModel "muslimly-be/internal/features/app_config/model"
	appConfigRepo "muslimly-be/internal/features/app_config/repository"
	appConfigService "muslimly-be/internal/features/app_config/service"

	articleHandler "muslimly-be/internal/features/article/handler"
	articleRepo "muslimly-be/internal/features/article/repository"
	articleService "muslimly-be/internal/features/article/service"

	"muslimly-be/pkg/config"
	"muslimly-be/pkg/database"
	"muslimly-be/pkg/logger"
	customMiddleware "muslimly-be/pkg/middleware"
)

// @title Muslimly Backend API
// @version 1.0
// @description Backend service for Muslimly App (Target Harian, Progress Sync, Notifications).
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@muslimly.app

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. LoadConfig
	cfg := config.LoadConfig()

	// 2. Initialize Logger
	logger.Init(logger.Config{
		EnableConsole: true,
		ConsoleJSON:   false,
		Verbose:       true,
	})
	log := logger.Log()
	log.Info().Msg("Starting Muslimly Backend API")

	// 3. Connect Database
	database.Connect(cfg)

	// 3. Auto Migrate
	if err := database.DB.AutoMigrate(
		&userModel.User{},
		&syncModel.ReadingHistory{},
		&syncModel.ReadingActivity{},
		&userSettingsModel.UserSettings{},
		&notifModel.UserDevice{},
		&appConfigModel.HijriAdjustment{},
		&articleRepo.Article{},
	); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	// 4. Setup Dependency Injection
	// User Feature
	uRepo := userRepo.NewUserRepository(database.DB)
	uService := userService.NewUserService(uRepo, cfg)
	uHandler := userHandler.NewUserHandler(uService)

	// Auth Feature
	aService := authService.NewAuthService(uRepo, cfg)
	aHandler := authHandler.NewAuthHandler(aService)

	// Sync Feature
	sRepo := syncRepo.NewSyncRepository(database.DB)
	sService := syncService.NewSyncService(sRepo, cfg)
	sHandler := syncHandler.NewSyncHandler(sService)

	// User Settings Feature
	usRepo := userSettingsRepo.NewUserSettingsRepository(database.DB)
	usService := userSettingsService.NewUserSettingsService(usRepo)
	usHandler := userSettingsHandler.NewUserSettingsHandler(usService)

	// Notification Feature
	nRepo := notifRepo.NewDeviceRepository(database.DB)
	nService := notifService.NewNotificationService(nRepo, cfg)
	nHandler := notifHandler.NewNotificationHandler(nService)

	// App Config Feature
	acRepo := appConfigRepo.NewAppConfigRepository(database.DB)
	acService := appConfigService.NewAppConfigService(cfg, acRepo)
	acHandler := appConfigHandler.NewAppConfigHandler(acService)

	// Article Feature
	artRepo := articleRepo.NewArticleRepository(database.DB)
	artService := articleService.NewArticleService(artRepo)
	artHandler := articleHandler.NewArticleHandler(artService)

	// 5. Initialize Echo
	e := echo.New()
	// e.Use(middleware.Logger()) // Replaced by RequestLogger
	e.Use(customMiddleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 6. Scheduler (Cron)
	c := cron.New()
	// Run everyday at 05:00 AM (0 5 * * *)
	// For testing, we can run every minute (* * * * *) if needed
	_, err := c.AddFunc("0 5 * * *", func() {
		log.Info().Msg("Executing Daily Reminder Job...")
		if err := nService.SendDailyReminder(); err != nil {
			log.Error().Err(err).Msg("Daily Reminder Failed")
		}
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to schedule daily reminder")
	}
	c.Start()
	log.Info().Msg("Scheduler started: Daily Reminder at 05:00 AM")

	// 7. Routes (Routes are defined in internal/api/router/router.go)

	// --- Rate Limiter Configuration Helpers ---
	rateLimitConfig := func(rateIn int) middleware.RateLimiterConfig {
		return middleware.RateLimiterConfig{
			Skipper: middleware.DefaultSkipper,
			Store: middleware.NewRateLimiterMemoryStoreWithConfig(
				middleware.RateLimiterMemoryStoreConfig{Rate: rate.Every(time.Minute / time.Duration(rateIn)), Burst: rateIn * 2, ExpiresIn: 3 * 60}, // 3 Minutes
			),
			IdentifierExtractor: func(ctx echo.Context) (string, error) {
				return ctx.RealIP(), nil
			},
			ErrorHandler: func(context echo.Context, err error) error {
				return context.JSON(429, map[string]string{"message": "Too many requests, please try again later."})
			},
			DenyHandler: func(context echo.Context, identifier string, err error) error {
				return context.JSON(429, map[string]string{"message": "Too many requests, please try again later."})
			},
		}
	}

	// Global Middleware (Default Limit)
	globalRate := cfg.RateLimit.Global
	if globalRate == 0 {
		globalRate = 60
	} // Default fallback

	config := rateLimitConfig(globalRate)
	config.Skipper = func(c echo.Context) bool {
		// Skip Global Limit for Public Articles (they have their own HIGHER limit)
		// We want to allow 100/min for articles, but global is 60/min.
		if len(c.Path()) > 0 && (c.Path() == "/api/v1/articles" || c.Request().URL.Path == "/api/v1/articles") {
			return true
		}
		// Also skip /swagger to avoid blocking docs
		if len(c.Path()) >= 8 && c.Path()[0:8] == "/swagger" {
			return true
		}
		return false
	}
	e.Use(middleware.RateLimiterWithConfig(config))

	appRouter := router.NewRouter(e, cfg)
	appRouter.RegisterRoutes(uHandler, aHandler, sHandler, usHandler, nHandler, acHandler, artHandler)

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// 8. Start Server (Gracerful Shutdown)
	// Start server in a goroutine so it doesn't block
	go func() {
		if err := e.Start(cfg.Server.Port); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop Scheduler
	c.Stop()
	log.Info().Msg("Scheduler stopped")

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}

func HealthCheck(c echo.Context) error {
	return c.JSON(200, map[string]string{"status": "ok"})
}
