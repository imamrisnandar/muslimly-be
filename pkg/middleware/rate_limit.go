package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

// RateLimitConfig creates a standard RateLimiter configuration
func RateLimitConfig(rateIn int) middleware.RateLimiterConfig {
	return middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Every(time.Minute / time.Duration(rateIn)),
				Burst:     rateIn * 2,
				ExpiresIn: 3 * time.Minute,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(429, map[string]string{
				"message": "Too many requests, please try again later.",
			})
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(429, map[string]string{
				"message": "Too many requests, please try again later.",
			})
		},
	}
}

// UserRateLimitConfig creates a RateLimiter configuration based on User ID
func UserRateLimitConfig(rateIn int) middleware.RateLimiterConfig {
	return middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      rate.Every(time.Minute / time.Duration(rateIn)),
				Burst:     rateIn * 2, // Allow burst of 2x the rate
				ExpiresIn: 3 * time.Minute,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			// Extract User ID from context (set by JWTMiddleware)
			userID, ok := ctx.Get("user_id").(string)
			if !ok || userID == "" {
				// Fallback to IP if User ID is missing (should not happen in protected routes)
				return ctx.RealIP(), nil
			}
			return userID, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(429, map[string]string{
				"message": "Too many sync requests. Please try again later.",
			})
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(429, map[string]string{
				"message": "Too many sync requests. Please try again later.",
			})
		},
	}
}
