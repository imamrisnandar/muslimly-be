package service

import (
	"muslimly-be/internal/features/user_settings/dto"
)

type UserSettingsService interface {
	UpsertSettings(userID, deviceID string, req dto.UpsertSettingsRequest) error
	GetSettings(userID, deviceID string) ([]dto.SettingResponse, error)
}
