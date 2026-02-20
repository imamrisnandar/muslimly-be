package middleware

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// RequestLogger middleware logs the HTTP request and response
func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// 1. Request ID (Generate if missing)
			reqID := c.Request().Header.Get(echo.HeaderXRequestID)
			if reqID == "" {
				reqID = uuid.New().String()
				c.Request().Header.Set(echo.HeaderXRequestID, reqID)
			}
			c.Response().Header().Set(echo.HeaderXRequestID, reqID)

			// 2. Attach Logger to Context with Request ID
			logger := log.With().Str("request_id", reqID).Logger()
			ctx := logger.WithContext(c.Request().Context())
			c.SetRequest(c.Request().WithContext(ctx))

			// 3. Process Request
			err := next(c)

			// 4. Log Response
			latency := time.Since(start)

			event := logger.Info()
			if err != nil {
				event = logger.Error().Err(err)
				c.Error(err) // Ensure error is handled by Echo
			}

			event.
				Str("method", c.Request().Method).
				Str("path", c.Request().URL.Path).
				Int("status", c.Response().Status).
				Dur("latency", latency).
				Str("ip", c.RealIP()).
				Str("user_agent", c.Request().UserAgent()).
				Msg("Request processed")

			return nil
		}
	}
}
