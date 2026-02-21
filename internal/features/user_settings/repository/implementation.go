package repository

import (
	"muslimly-be/internal/features/user_settings/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userSettingsRepository struct {
	db *gorm.DB
}

func NewUserSettingsRepository(db *gorm.DB) UserSettingsRepository {
	return &userSettingsRepository{db}
}

func (r *userSettingsRepository) UpsertSettings(settings []model.UserSettings) error {
	// Batch Upsert using the new composite unique index
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "device_id"}, {Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&settings).Error
}

func (r *userSettingsRepository) GetSettings(userID, deviceID string) ([]model.UserSettings, error) {
	var settings []model.UserSettings

	query := r.db
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	} else if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	} else {
		// Prevent accidental fetch of everything if both are empty
		return nil, gorm.ErrRecordNotFound
	}

	err := query.Find(&settings).Error
	return settings, err
}
