package handler

import (
	"net/http"

	"muslimly-be/internal/features/notification/dto"
	"muslimly-be/internal/features/notification/service"
	"muslimly-be/pkg/utils"

	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	service service.NotificationService
}

func NewNotificationHandler(service service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service}
}

// RegisterDevice godoc
// @Summary Register FCM Device Token
// @Description Register device for push notifications
// @Tags Notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.RegisterDeviceRequest true "Device Info"
// @Success 200 {object} utils.WebResponse
// @Router /notifications/register [post]
func (h *NotificationHandler) RegisterDevice(c echo.Context) error {
	// Optional Auth: GetUserIDFromContext might return empty if no token
	userID := utils.GetUserIDFromContext(c)

	var req dto.RegisterDeviceRequest
	if err := c.Bind(&req); err != nil {
		return utils.ResponseError(c, http.StatusBadRequest, utils.ErrInvalidRequest, nil)
	}

	if err := h.service.RegisterDevice(userID, req); err != nil {
		return utils.ResponseError(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseSuccess(c, http.StatusOK, "Device registered", nil)
}

// TestBroadcast godoc
// @Summary Trigger Daily Reminder (Debug)
// @Description Force send daily reminder to all devices
// @Tags Notification
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.WebResponse
// @Router /notifications/test-broadcast [post]
func (h *NotificationHandler) TestBroadcast(c echo.Context) error {
	if err := h.service.SendDailyReminder(); err != nil {
		return utils.ResponseError(c, http.StatusInternalServerError, err.Error(), nil)
	}
	return utils.ResponseSuccess(c, http.StatusOK, "Broadcast triggered", nil)
}
