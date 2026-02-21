package repository

import (
	"muslimly-be/internal/features/sync/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type syncRepository struct {
	db *gorm.DB
}

func NewSyncRepository(db *gorm.DB) SyncRepository {
	return &syncRepository{db}
}

func (r *syncRepository) UpsertReading(history *model.ReadingHistory) error {
	history.LastReadAt = time.Now()

	// Determine conflict columns based on identity
	if history.UserID != nil {
		// Logged-in user: upsert based on (user_id, surah_id)
		return r.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "surah_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"ayah_number", "page_number", "mode", "last_read_at", "updated_at", "device_id"}),
		}).Create(history).Error
	} else if history.DeviceID != nil {
		// Guest: upsert based on (device_id, surah_id)
		return r.db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "device_id"}, {Name: "surah_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"ayah_number", "page_number", "mode", "last_read_at", "updated_at"}),
		}).Create(history).Error
	}

	// Fallback: just insert
	return r.db.Create(history).Error
}

func (r *syncRepository) GetReadingHistory(userID string, deviceID string, limit int) ([]model.ReadingHistory, error) {
	var histories []model.ReadingHistory
	query := r.db.Order("last_read_at desc").Limit(limit)

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	} else if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	} else {
		return histories, nil // No identity, return empty
	}

	err := query.Find(&histories).Error
	return histories, err
}

func (r *syncRepository) UpsertActivities(activities []model.ReadingActivity) error {
	return r.db.Create(&activities).Error
}
