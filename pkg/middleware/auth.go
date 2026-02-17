package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"muslimly-be/pkg/config"
	"muslimly-be/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return utils.ResponseError(c, http.StatusUnauthorized, "Missing Authorization header", nil)
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return utils.ResponseError(c, http.StatusUnauthorized, "Invalid token format", nil)
			}

			tokenString := parts[1]
			token, err := jwt.ParseWithClaims(tokenString, &utils.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(cfg.JWT.Secret), nil
			})

			if err != nil || !token.Valid {
				return utils.ResponseError(c, http.StatusUnauthorized, "Invalid or expired token", nil)
			}

			claims, ok := token.Claims.(*utils.JWTClaims)
			if !ok {
				return utils.ResponseError(c, http.StatusUnauthorized, "Invalid token claims", nil)
			}

			// Store user identity in context
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}

// OptionalJWTMiddleware attempts to validate JWT but proceeds even if invalid/missing
func OptionalJWTMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				// No token, proceed as Guest (userID context will be empty)
				return next(c)
			}

			// Validate format "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				// Invalid format, but since optional, we treat as Guest
				return next(c)
			}

			tokenString := parts[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.ErrUnauthorized
				}
				return []byte(cfg.JWT.Secret), nil
			})

			if err != nil || !token.Valid {
				// Invalid token, treat as Guest
				return next(c)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return next(c)
			}

			userID, ok := claims["user_id"].(string)
			if !ok {
				return next(c)
			}

			// Set userID in context if valid
			c.Set("user_id", userID)
			return next(c)
		}
	}
}
