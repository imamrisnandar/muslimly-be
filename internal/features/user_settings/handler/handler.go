package handler

import (
	"net/http"

	"muslimly-be/internal/features/user_settings/dto"
	"muslimly-be/internal/features/user_settings/service"
	"muslimly-be/pkg/utils"

	"github.com/labstack/echo/v4"
)

type UserSettingsHandler struct {
	service service.UserSettingsService
}

func NewUserSettingsHandler(service service.UserSettingsService) *UserSettingsHandler {
	return &UserSettingsHandler{service}
}

// UpsertSettings godoc
// @Summary Upsert User Settings
// @Description Update user settings (key-value pairs)
// @Tags Sync
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpsertSettingsRequest true "Settings"
// @Success 200 {object} utils.WebResponse
// @Router /sync/settings [post]
func (h *UserSettingsHandler) UpsertSettings(c echo.Context) error {
	userID := utils.GetUserIDFromContext(c)
	var req dto.UpsertSettingsRequest
	if err := c.Bind(&req); err != nil {
		return utils.ResponseError(c, http.StatusBadRequest, utils.ErrInvalidRequest, nil)
	}

	if err := h.service.UpsertSettings(userID, req); err != nil {
		return utils.ResponseError(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseSuccess(c, http.StatusOK, "Settings updated", nil)
}

// GetSettings godoc
// @Summary Get User Settings
// @Description Get all user settings
// @Tags Sync
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.WebResponse{data=[]dto.SettingResponse}
// @Router /sync/settings [get]
func (h *UserSettingsHandler) GetSettings(c echo.Context) error {
	userID := utils.GetUserIDFromContext(c)

	settings, err := h.service.GetSettings(userID)
	if err != nil {
		return utils.ResponseError(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseSuccess(c, http.StatusOK, "Settings retrieved", settings)
}
