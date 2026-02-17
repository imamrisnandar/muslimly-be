package service

import (
	authdto "muslimly-be/internal/features/auth/dto"
	"muslimly-be/internal/features/user/model"
)

type AuthService interface {
	Register(req authdto.RegisterRequest) (*model.User, error)
	Login(req authdto.LoginRequest) (string, *model.User, error)
}
