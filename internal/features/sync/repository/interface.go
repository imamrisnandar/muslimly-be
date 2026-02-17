package repository

import "muslimly-be/internal/features/sync/model"

type SyncRepository interface {
	UpsertReading(history *model.ReadingHistory) error
	GetReadingHistory(userID string, limit int) ([]model.ReadingHistory, error)
	UpsertActivities(activities []model.ReadingActivity) error
}
