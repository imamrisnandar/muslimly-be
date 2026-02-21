package router

import (
	appConfigHandler "muslimly-be/internal/features/app_config/handler"
	"muslimly-be/internal/features/article/handler"
	authHandler "muslimly-be/internal/features/auth/handler"
	notifHandler "muslimly-be/internal/features/notification/handler"
	syncHandler "muslimly-be/internal/features/sync/handler"
	userHandler "muslimly-be/internal/features/user/handler"
	userSettingsHandler "muslimly-be/internal/features/user_settings/handler"
	"muslimly-be/pkg/config"
	customMiddleware "muslimly-be/pkg/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	articleHandler *handler.ArticleHandler,
) {
	v1 := r.echo.Group("/api/v1")

	// Health Check
	v1.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "up",
			"message": "Muslimly Backend is running 🚀",
		})
	})

	// ... (Auth Routes) ...
	// Public Routes (Auth) - Limit: 10/min (or config)
	authRate := r.config.RateLimit.Auth
	if authRate == 0 {
		authRate = 10
	}
	auth := v1.Group("/auth")
	auth.Use(middleware.RateLimiterWithConfig(customMiddleware.RateLimitConfig(authRate)))
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

	// Sync Routes (Supports both logged-in users and guests via device_id)
	sync := v1.Group("/sync")
	sync.Use(customMiddleware.OptionalJWTMiddleware(r.config))

	// Sync Limit (User Based)
	syncRate := r.config.RateLimit.Sync
	if syncRate == 0 {
		syncRate = 20
	}
	sync.Use(middleware.RateLimiterWithConfig(customMiddleware.UserRateLimitConfig(syncRate)))

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

	// Article (Public) - Limit: 100/min (or config)
	publicRate := r.config.RateLimit.Public
	if publicRate == 0 {
		publicRate = 100
	}
	v1.GET("/articles", articleHandler.GetArticles, middleware.RateLimiterWithConfig(customMiddleware.RateLimitConfig(publicRate)))
}
