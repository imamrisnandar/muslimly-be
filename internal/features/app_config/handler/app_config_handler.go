package handler

import (
	"net/http"

	"muslimly-be/internal/features/app_config/service"

	"github.com/labstack/echo/v4"
)

type AppConfigHandler struct {
	service service.AppConfigService
}

func NewAppConfigHandler(service service.AppConfigService) *AppConfigHandler {
	return &AppConfigHandler{
		service: service,
	}
}

// GetAppConfig godoc
// @Summary      Get Hijri adjustment configuration
// @Description  Get a list of Hijri adjustments (year-agnostic, by month)
// @Tags         App Config
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.AppConfigResponse
// @Router       /config-hijri-adjust [get]
func (h *AppConfigHandler) GetAppConfig(c echo.Context) error {
	config, err := h.service.GetAppConfig()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, config)
}
