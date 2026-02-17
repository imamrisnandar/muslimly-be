package service

import (
	"errors"
	"muslimly-be/internal/features/user/dto"
	"muslimly-be/internal/features/user/model"
	"muslimly-be/internal/features/user/repository"
	"muslimly-be/pkg/config"
	"muslimly-be/pkg/utils"
)

type userService struct {
	repo   repository.UserRepository
	config *config.Config
}

func NewUserService(repo repository.UserRepository, config *config.Config) UserService {
	return &userService{repo, config}
}

func (s *userService) Update(req dto.UpdateUserRequest, actorID string) (*model.User, error) {
	user, err := s.repo.FindByID(req.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if req.Username != "" {
		user.Username = req.Username
	}

	if req.Email != "" && req.Email != user.Email {
		// Check uniqueness
		existing, _ := s.repo.FindByEmail(req.Email)
		if existing != nil {
			return nil, errors.New(utils.ErrEmailExists)
		}
		user.Email = req.Email
	}

	// Audit Trail
	user.UpdatedBy = actorID

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(id string, actorID string) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	return s.repo.Delete(id, actorID)
}

func (s *userService) GetByID(id string) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) GetAll(req dto.GetDataRequest) (*dto.ListUserResponse, error) {
	// Defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	users, total, err := s.repo.FindAll(req.Page, req.Limit, req.Sort, req.Filters)
	if err != nil {
		return nil, err
	}

	// Map to DTO
	var list []dto.UserResponse
	for _, u := range users {
		list = append(list, dto.UserResponse{
			ID:       u.ID.String(),
			Username: u.Username,
			Email:    u.Email,
		})
	}

	// Calculate Total Page
	totalPage := int(total) / req.Limit
	if int(total)%req.Limit != 0 {
		totalPage++
	}

	return &dto.ListUserResponse{
		List: list,
		Meta: dto.PaginationMeta{
			CurrentPage: req.Page,
			TotalPage:   totalPage,
			TotalData:   total,
			Limit:       req.Limit,
		},
	}, nil
}
