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

func (s *userSettingsService) UpsertSettings(userID, deviceID string, req dto.UpsertSettingsRequest) error {
	var uidPtr *uuid.UUID
	if userID != "" {
		uid, err := uuid.Parse(userID)
		if err == nil {
			uidPtr = &uid
		}
	}

	var parsedDeviceID *uuid.UUID
	if deviceID != "" {
		did, err := uuid.Parse(deviceID)
		if err == nil {
			parsedDeviceID = &did
		}
	}

	var settings []model.UserSettings
	for _, item := range req.Settings {
		settings = append(settings, model.UserSettings{
			UserID:   uidPtr,
			DeviceID: parsedDeviceID,
			Key:      item.Key,
			Value:    item.Value,
		})
	}

	return s.repo.UpsertSettings(settings)
}

func (s *userSettingsService) GetSettings(userID, deviceID string) ([]dto.SettingResponse, error) {
	data, err := s.repo.GetSettings(userID, deviceID)
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
