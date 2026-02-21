package repository

import (
	"muslimly-be/internal/features/user_settings/model"
)

type UserSettingsRepository interface {
	UpsertSettings(settings []model.UserSettings) error
	GetSettings(userID, deviceID string) ([]model.UserSettings, error)
}
