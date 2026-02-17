package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
	echoSwagger "github.com/swaggo/echo-swagger"

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

	"muslimly-be/pkg/config"
	"muslimly-be/pkg/database"
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
	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Connect Database
	database.Connect(cfg)

	// 3. Auto Migrate
	if err := database.DB.AutoMigrate(
		&userModel.User{},
		&syncModel.ReadingHistory{},
		&syncModel.ReadingActivity{},
		&userSettingsModel.UserSettings{},
		&notifModel.UserDevice{},
		&appConfigModel.HijriAdjustment{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
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

	// 5. Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// 6. Scheduler (Cron)
	c := cron.New()
	// Run everyday at 05:00 AM (0 5 * * *)
	// For testing, we can run every minute (* * * * *) if needed
	_, err := c.AddFunc("0 5 * * *", func() {
		log.Println("Executing Daily Reminder Job...")
		if err := nService.SendDailyReminder(); err != nil {
			log.Printf("Daily Reminder Failed: %v", err)
		}
	})
	if err != nil {
		log.Printf("Failed to schedule daily reminder: %v", err)
	}
	c.Start()
	log.Println("Scheduler started: Daily Reminder at 05:00 AM")

	// 7. Routes
	v1 := e.Group("/api/v1")
	{
		v1.GET("/health", HealthCheck)
	}
	appRouter := router.NewRouter(e, cfg)
	appRouter.RegisterRoutes(uHandler, aHandler, sHandler, usHandler, nHandler, acHandler)

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// 7. Start Server
	e.Logger.Fatal(e.Start(cfg.Server.Port))
}

func HealthCheck(c echo.Context) error {
	return c.JSON(200, map[string]string{"status": "ok"})
}
