package utils

import "github.com/labstack/echo/v4"

func GetUserIDFromContext(c echo.Context) string {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return ""
	}
	return userID
}

func GetUserEmailFromContext(c echo.Context) string {
	email, ok := c.Get("email").(string)
	if !ok {
		return ""
	}
	return email
}
