package dto

import "time"

type UpsertReadingRequest struct {
	DeviceID   string `json:"device_id"` // Device UUID for guest sync
	SurahID    int    `json:"surah_id" validate:"required"`
	AyahNumber int    `json:"ayah_number" validate:"required"`
	PageNumber int    `json:"page_number"`
	Mode       string `json:"mode"` // 'mushaf', 'list'
}

type ReadingActivityRequest struct {
	Date            string `json:"date" validate:"required"` // YYYY-MM-DD
	DurationSeconds int    `json:"duration_seconds"`
	PageNumber      int    `json:"page_number"`
	SurahNumber     int    `json:"surah_number"`
	StartAyah       int    `json:"start_ayah"`
	EndAyah         int    `json:"end_ayah"`
	TotalAyahs      int    `json:"total_ayahs"`
	Mode            string `json:"mode"`
	Timestamp       int64  `json:"timestamp"`
}

type BulkActivityRequest struct {
	DeviceID   string                   `json:"device_id"` // Device UUID for guest sync
	Activities []ReadingActivityRequest `json:"activities" validate:"required,dive"`
}

type ReadingHistoryResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id,omitempty"`
	DeviceID   string    `json:"device_id,omitempty"`
	SurahID    int       `json:"surah_id"`
	AyahNumber int       `json:"ayah_number"`
	PageNumber int       `json:"page_number"`
	Mode       string    `json:"mode"`
	LastReadAt time.Time `json:"last_read_at"`
}
