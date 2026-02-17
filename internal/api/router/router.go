package router

import (
	appConfigHandler "muslimly-be/internal/features/app_config/handler"
	authHandler "muslimly-be/internal/features/auth/handler"
	notifHandler "muslimly-be/internal/features/notification/handler"
	syncHandler "muslimly-be/internal/features/sync/handler"
	userHandler "muslimly-be/internal/features/user/handler"
	userSettingsHandler "muslimly-be/internal/features/user_settings/handler"
	"muslimly-be/pkg/config"
	customMiddleware "muslimly-be/pkg/middleware"

	"github.com/labstack/echo/v4"
)

type Router struct {
	echo   *echo.Echo
	config *config.Config
}

func NewRouter(echo *echo.Echo, config *config.Config) *Router {
	return &Router{echo, config}
}

func (r *Router) RegisterRoutes(
	userHandler *userHandler.UserHandler,
	authHandler *authHandler.AuthHandler,
	syncHandler *syncHandler.SyncHandler,
	userSettingsHandler *userSettingsHandler.UserSettingsHandler,
	notifHandler *notifHandler.NotificationHandler,
	appConfigHandler *appConfigHandler.AppConfigHandler,
) {
	v1 := r.echo.Group("/api/v1")

	// Health Check
	v1.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "up",
			"message": "Muslimly Backend is running 🚀",
		})
	})

	// Public Routes (Auth)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected Routes (User)
	users := v1.Group("/users")
	users.Use(customMiddleware.JWTMiddleware(r.config))
	{
		users.PUT("/update", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
		users.GET("/:id", userHandler.GetByID)
		users.POST("/list", userHandler.GetData)
	}

	// Protected Routes (Sync)
	sync := v1.Group("/sync")
	sync.Use(customMiddleware.JWTMiddleware(r.config))
	{
		// Reading History
		sync.POST("/reading", syncHandler.UpsertReading)
		sync.GET("/reading", syncHandler.GetReadingHistory)
		sync.POST("/activity", syncHandler.BulkInsertActivities)

		// User Settings
		sync.POST("/settings", userSettingsHandler.UpsertSettings)
		sync.GET("/settings", userSettingsHandler.GetSettings)
	}

	// Protected Routes (Notification) - Mixed (Public/Protected)
	notif := v1.Group("/notifications")
	// REGISTER IS PUBLIC now to support Guest
	notif.POST("/register", notifHandler.RegisterDevice, customMiddleware.OptionalJWTMiddleware(r.config))

	// Protected Test endpoint
	notif.POST("/test-broadcast", notifHandler.TestBroadcast, customMiddleware.JWTMiddleware(r.config))

	// App Config (Public)
	v1.GET("/config-hijri-adjust", appConfigHandler.GetAppConfig)
}
