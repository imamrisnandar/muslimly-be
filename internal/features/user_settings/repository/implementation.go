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
	// Batch Upsert
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&settings).Error
}

func (r *userSettingsRepository) GetSettings(userID string) ([]model.UserSettings, error) {
	var settings []model.UserSettings
	err := r.db.Where("user_id = ?", userID).Find(&settings).Error
	return settings, err
}
