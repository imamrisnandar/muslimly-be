package service

import (
	"muslimly-be/internal/features/user/dto"
	"muslimly-be/internal/features/user/model"
)

type UserService interface {
	Update(req dto.UpdateUserRequest, actorID string) (*model.User, error)
	Delete(id string, actorID string) error
	GetByID(id string) (*model.User, error)
	GetAll(req dto.GetDataRequest) (*dto.ListUserResponse, error)
}
