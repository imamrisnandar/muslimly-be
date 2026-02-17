package service

import (
	"muslimly-be/internal/features/user_settings/dto"
	"muslimly-be/internal/features/user_settings/model"
	"muslimly-be/internal/features/user_settings/repository"

	"github.com/google/uuid"
)

type userSettingsService struct {
	repo repository.UserSettingsRepository
}

func NewUserSettingsService(repo repository.UserSettingsRepository) UserSettingsService {
	return &userSettingsService{repo}
}

func (s *userSettingsService) UpsertSettings(userID string, req dto.UpsertSettingsRequest) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	var settings []model.UserSettings
	for _, item := range req.Settings {
		settings = append(settings, model.UserSettings{
			UserID: uid,
			Key:    item.Key,
			Value:  item.Value,
		})
	}

	return s.repo.UpsertSettings(settings)
}

func (s *userSettingsService) GetSettings(userID string) ([]dto.SettingResponse, error) {
	data, err := s.repo.GetSettings(userID)
	if err != nil {
		return nil, err
	}

	var response []dto.SettingResponse
	for _, item := range data {
		response = append(response, dto.SettingResponse{
			Key:   item.Key,
			Value: item.Value,
		})
	}
	return response, nil
}
