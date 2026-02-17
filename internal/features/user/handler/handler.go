package handler

import (
	"net/http"

	"muslimly-be/internal/features/user/dto"
	"muslimly-be/internal/features/user/service"
	"muslimly-be/pkg/utils"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service}
}

// Update godoc
// @Summary Update user
// @Description Update user details
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.UpdateUserRequest true "Update Request"
// @Success 200 {object} utils.WebResponse{data=dto.UserResponse}
// @Failure 400 {object} utils.WebResponse
// @Security BearerAuth
// @Router /users/update [put]
func (h *UserHandler) Update(c echo.Context) error {
	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return utils.ResponseError(c, http.StatusBadRequest, utils.ErrInvalidRequest, nil)
	}

	// Get Actor ID from Context (Set by Middleware)
	actorID := utils.GetUserIDFromContext(c)

	user, err := h.service.Update(req, actorID)
	if err != nil {
		return utils.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.ResponseSuccess(c, http.StatusOK, "User updated successfully", dto.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	})
}

// Delete godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} utils.WebResponse
// @Failure 400 {object} utils.WebResponse
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.ResponseError(c, http.StatusBadRequest, "ID is required", nil)
	}

	// Get Actor ID form Context
	actorID := utils.GetUserIDFromContext(c)

	if err := h.service.Delete(id, actorID); err != nil {
		return utils.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.ResponseSuccess(c, http.StatusOK, "User deleted successfully", nil)
}

// GetByID godoc
// @Summary Get user by ID
// @Description Get single user details
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} utils.WebResponse{data=dto.UserResponse}
// @Failure 404 {object} utils.WebResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	user, err := h.service.GetByID(id)
	if err != nil {
		return utils.ResponseError(c, http.StatusNotFound, "User not found", nil)
	}

	return utils.ResponseSuccess(c, http.StatusOK, "User retrieved successfully", dto.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	})
}

// GetData godoc
// @Summary List users
// @Description Get users with pagination, filter, and sort via Body
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.GetDataRequest true "Filter Request"
// @Success 200 {object} utils.WebResponse{data=dto.ListUserResponse}
// @Failure 400 {object} utils.WebResponse
// @Security BearerAuth
// @Router /users/list [post]
func (h *UserHandler) GetData(c echo.Context) error {
	var req dto.GetDataRequest
	if err := c.Bind(&req); err != nil {
		return utils.ResponseError(c, http.StatusBadRequest, utils.ErrInvalidRequest, nil)
	}

	resp, err := h.service.GetAll(req)
	if err != nil {
		return utils.ResponseError(c, http.StatusInternalServerError, utils.ErrInternalServer, nil)
	}

	return utils.ResponseSuccess(c, http.StatusOK, "Data retrieved successfully", resp)
}
