package handler

import (
	"net/http"

	authdto "muslimly-be/internal/features/auth/dto"
	"muslimly-be/internal/features/auth/service"
	"muslimly-be/pkg/utils"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authdto.RegisterRequest true "Register Request"
// @Success 201 {object} utils.WebResponse{data=authdto.UserResponse}
// @Failure 400 {object} utils.WebResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req authdto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.ResponseError(c, http.StatusBadRequest, utils.ErrInvalidRequest, nil)
	}

	user, err := h.service.Register(req)
	if err != nil {
		return utils.ResponseError(c, http.StatusBadRequest, err.Error(), nil)
	}

	return utils.ResponseSuccess(c, http.StatusCreated, utils.MsgUserRegistered, authdto.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	})
}

// Login godoc
// @Summary Login user
// @Description Login with email and password to get JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authdto.LoginRequest true "Login Request"
// @Success 200 {object} utils.WebResponse{data=authdto.AuthResponse}
// @Failure 401 {object} utils.WebResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req authdto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.ResponseError(c, http.StatusBadRequest, utils.ErrInvalidRequest, nil)
	}

	token, user, err := h.service.Login(req)
	if err != nil {
		return utils.ResponseError(c, http.StatusUnauthorized, err.Error(), nil)
	}

	return utils.ResponseSuccess(c, http.StatusOK, utils.MsgLoginSuccess, authdto.AuthResponse{
		Token: token,
		User: authdto.UserResponse{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
		},
	})
}
