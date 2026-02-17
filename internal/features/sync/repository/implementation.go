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
	// Use GORM Clauses for Upsert
	// Requires unique index on (user_id, surah_id)
	history.LastReadAt = time.Now()

	// If conflict on (user_id, surah_id), update all progress fields
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "surah_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"ayah_number", "page_number", "mode", "last_read_at", "updated_at"}),
	}).Create(history).Error
}

func (r *syncRepository) GetReadingHistory(userID string, limit int) ([]model.ReadingHistory, error) {
	var histories []model.ReadingHistory
	err := r.db.Where("user_id = ?", userID).
		Order("last_read_at desc").
		Limit(limit).
		Find(&histories).Error
	return histories, err
}

func (r *syncRepository) UpsertActivities(activities []model.ReadingActivity) error {
	// Batch Insert. We assume these are logs.
	// If mobile sends same log twice, we might want to avoid dupes?
	// Mobile has 'id' but it is local sqlite ID.
	// We can use (user_id, timestamp) as unique key? Or just append?
	// For now, let's just append or maybe use simple unique constraint on timestamp if critical.
	// But users might read multiple times in same second (unlikely but possible via batch sync).
	// Let's rely on client not sending duplicates or standard ID generation.
	// Actually, best practice is to have a unique ID or composite key.
	// Let's assume we just insert for now. Mobile DB has ID, we could map it if we want strict sync.

	return r.db.Create(&activities).Error
}
