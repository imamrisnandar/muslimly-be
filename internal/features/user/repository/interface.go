package repository

import "muslimly-be/internal/features/user/model"

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByID(id string) (*model.User, error)
	Update(user *model.User) error
	Delete(id string, actorID string) error
	FindAll(page, limit int, sort string, filters map[string]interface{}) ([]model.User, int64, error)
}
