package service

import (
	"muslimly-be/internal/features/sync/dto"
)

type SyncService interface {
	UpsertReading(userID string, deviceID string, req dto.UpsertReadingRequest) error
	GetReadingHistory(userID string, deviceID string) ([]dto.ReadingHistoryResponse, error)
	BulkInsertActivities(userID string, deviceID string, req dto.BulkActivityRequest) error
}
