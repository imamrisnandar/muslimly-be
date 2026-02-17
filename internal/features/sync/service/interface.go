package service

import (
	"muslimly-be/internal/features/sync/dto"
)

type SyncService interface {
	UpsertReading(userID string, req dto.UpsertReadingRequest) error
	GetReadingHistory(userID string) ([]dto.ReadingHistoryResponse, error)
	BulkInsertActivities(userID string, req dto.BulkActivityRequest) error
}
