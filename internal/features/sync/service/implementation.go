package service

import (
	"errors"
	"muslimly-be/internal/features/sync/dto"
	"muslimly-be/internal/features/sync/model"
	"muslimly-be/internal/features/sync/repository"
	"muslimly-be/pkg/config"

	"github.com/google/uuid"
)

type syncService struct {
	repo   repository.SyncRepository
	config *config.Config
}

func NewSyncService(repo repository.SyncRepository, config *config.Config) SyncService {
	return &syncService{repo, config}
}

// resolveIdentity parses userID and deviceID strings into UUID pointers.
// At least one must be valid. Returns (userUUID, deviceUUID, error).
func resolveIdentity(userID, deviceID string) (*uuid.UUID, *uuid.UUID, error) {
	var uid *uuid.UUID
	var did *uuid.UUID

	if userID != "" {
		parsed, err := uuid.Parse(userID)
		if err == nil {
			uid = &parsed
		}
	}

	if deviceID != "" {
		parsed, err := uuid.Parse(deviceID)
		if err == nil {
			did = &parsed
		}
	}

	if uid == nil && did == nil {
		return nil, nil, errors.New("either user_id or device_id is required")
	}

	return uid, did, nil
}

func (s *syncService) UpsertReading(userID string, deviceID string, req dto.UpsertReadingRequest) error {
	uid, did, err := resolveIdentity(userID, deviceID)
	if err != nil {
		return err
	}

	history := &model.ReadingHistory{
		UserID:     uid,
		DeviceID:   did,
		SurahID:    req.SurahID,
		AyahNumber: req.AyahNumber,
		PageNumber: req.PageNumber,
		Mode:       req.Mode,
	}

	// Default mode if empty
	if history.Mode == "" {
		history.Mode = "mushaf"
	}

	return s.repo.UpsertReading(history)
}

func (s *syncService) GetReadingHistory(userID string, deviceID string) ([]dto.ReadingHistoryResponse, error) {
	histories, err := s.repo.GetReadingHistory(userID, deviceID, 10) // Limit 10 recent
	if err != nil {
		return nil, err
	}

	var responses []dto.ReadingHistoryResponse
	for _, h := range histories {
		resp := dto.ReadingHistoryResponse{
			ID:         h.ID.String(),
			SurahID:    h.SurahID,
			AyahNumber: h.AyahNumber,
			PageNumber: h.PageNumber,
			Mode:       h.Mode,
			LastReadAt: h.LastReadAt,
		}
		if h.UserID != nil {
			resp.UserID = h.UserID.String()
		}
		if h.DeviceID != nil {
			resp.DeviceID = h.DeviceID.String()
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

func (s *syncService) BulkInsertActivities(userID string, deviceID string, req dto.BulkActivityRequest) error {
	uid, did, err := resolveIdentity(userID, deviceID)
	if err != nil {
		return err
	}

	var activities []model.ReadingActivity
	for _, item := range req.Activities {
		activities = append(activities, model.ReadingActivity{
			UserID:          uid,
			DeviceID:        did,
			Date:            item.Date,
			DurationSeconds: item.DurationSeconds,
			PageNumber:      item.PageNumber,
			SurahNumber:     item.SurahNumber,
			StartAyah:       item.StartAyah,
			EndAyah:         item.EndAyah,
			TotalAyahs:      item.TotalAyahs,
			Mode:            item.Mode,
			Timestamp:       item.Timestamp,
		})
	}

	if len(activities) == 0 {
		return nil
	}

	return s.repo.UpsertActivities(activities)
}
