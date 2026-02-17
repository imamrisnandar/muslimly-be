package handler

import (
	"net/http"

	"muslimly-be/internal/features/sync/dto"
	"muslimly-be/internal/features/sync/service"
	"muslimly-be/pkg/utils"

	"github.com/labstack/echo/v4"
)

type SyncHandler struct {
	service service.SyncService
}

func NewSyncHandler(service service.SyncService) *SyncHandler {
	return &SyncHandler{service}
}

// UpsertReading godoc
// @Summary Sync Reading History
// @Description Upsert reading history (Last Read Position)
// @Tags Sync
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpsertReadingRequest true "Reading History"
// @Success 200 {object} utils.WebResponse
// @Router /sync/reading [post]
func (h *SyncHandler) UpsertReading(c echo.Context) error {
	userID := utils.GetUserIDFromContext(c)
	var req dto.UpsertReadingRequest
	if err := c.Bind(&req); err != nil {
		return utils.ResponseError(c, http.StatusBadRequest, utils.ErrInvalidRequest, nil)
	}

	if err := h.service.UpsertReading(userID, req); err != nil {
		return utils.ResponseError(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseSuccess(c, http.StatusOK, "Reading progress updated", nil)
}

// GetReadingHistory godoc
// @Summary Get Reading History
// @Description Get last 10 reading history items
// @Tags Sync
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.WebResponse{data=[]dto.ReadingHistoryResponse}
// @Router /sync/reading [get]
func (h *SyncHandler) GetReadingHistory(c echo.Context) error {
	userID := utils.GetUserIDFromContext(c)

	history, err := h.service.GetReadingHistory(userID)
	if err != nil {
		return utils.ResponseError(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseSuccess(c, http.StatusOK, "Reading history retrieved", history)
}

// BulkInsertActivities godoc
// @Summary Bulk Sync Reading Activities
// @Description Upload multiple reading activity logs
// @Tags Sync
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.BulkActivityRequest true "Activities"
// @Success 200 {object} utils.WebResponse
// @Router /sync/activity [post]
func (h *SyncHandler) BulkInsertActivities(c echo.Context) error {
	userID := utils.GetUserIDFromContext(c)
	var req dto.BulkActivityRequest
	if err := c.Bind(&req); err != nil {
		return utils.ResponseError(c, http.StatusBadRequest, utils.ErrInvalidRequest, nil)
	}

	if err := h.service.BulkInsertActivities(userID, req); err != nil {
		return utils.ResponseError(c, http.StatusInternalServerError, err.Error(), nil)
	}

	return utils.ResponseSuccess(c, http.StatusOK, "Activities synced", nil)
}
