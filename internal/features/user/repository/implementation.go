package repository

import (
	"errors"
	"muslimly-be/internal/features/user/model"
	"muslimly-be/pkg/utils"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id string, actorID string) error {
	// Soft Delete with Audit
	// 1. Update DeletedBy
	if err := r.db.Model(&model.User{}).Where("id = ?", id).Update("deleted_by", actorID).Error; err != nil {
		return err
	}
	// 2. Perform Soft Delete (GORM handles DeletedAt automatically)
	return r.db.Where("id = ?", id).Delete(&model.User{}).Error
}

func (r *userRepository) FindAll(page, limit int, sort string, filters map[string]interface{}) ([]model.User, int64, error) {
	var users []model.User

	query := r.db.Model(&model.User{})

	// Use generic pagination helper
	query, total := utils.Paginate(query, utils.PaginationConfig{
		Page:    page,
		Limit:   limit,
		Sort:    sort,
		Filters: filters,
	})

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
