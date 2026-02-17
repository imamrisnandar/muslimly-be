package service

import (
	"errors"
	authdto "muslimly-be/internal/features/auth/dto"
	"muslimly-be/internal/features/user/model"
	"muslimly-be/internal/features/user/repository"
	"muslimly-be/pkg/config"
	"muslimly-be/pkg/utils"
)

type authService struct {
	userRepo repository.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo repository.UserRepository, config *config.Config) AuthService {
	return &authService{userRepo, config}
}

func (s *authService) Register(req authdto.RegisterRequest) (*model.User, error) {
	// Check if user exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New(utils.ErrEmailExists)
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(req authdto.LoginRequest) (string, *model.User, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New(utils.ErrInvalidCreds)
	}

	if !user.CheckPassword(req.Password) {
		return "", nil, errors.New(utils.ErrInvalidCreds)
	}

	// Generate Token
	token, err := utils.GenerateToken(user.ID.String(), user.Email, s.config.JWT.Secret, s.config.JWT.ExpirationHours)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
