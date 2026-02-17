package service

import (
	"muslimly-be/internal/features/user_settings/dto"
)

type UserSettingsService interface {
	UpsertSettings(userID string, req dto.UpsertSettingsRequest) error
	GetSettings(userID string) ([]dto.SettingResponse, error)
}
