package utils

import "github.com/labstack/echo/v4"

type WebResponse struct {
	Code        int         `json:"code"`
	Status      bool        `json:"status"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data,omitempty"`
	ErrorFields interface{} `json:"error_fields,omitempty"`
}

func ResponseSuccess(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, WebResponse{
		Code:    code,
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func ResponseError(c echo.Context, code int, message string, errorFields interface{}) error {
	return c.JSON(code, WebResponse{
		Code:        code,
		Status:      false,
		Message:     message,
		ErrorFields: errorFields,
	})
}
