package service

import (
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

func (s *syncService) UpsertReading(userID string, req dto.UpsertReadingRequest) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	history := &model.ReadingHistory{
		UserID:     uid,
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

func (s *syncService) GetReadingHistory(userID string) ([]dto.ReadingHistoryResponse, error) {
	histories, err := s.repo.GetReadingHistory(userID, 10) // Limit 10 recent
	if err != nil {
		return nil, err
	}

	var responses []dto.ReadingHistoryResponse
	for _, h := range histories {
		responses = append(responses, dto.ReadingHistoryResponse{
			ID:         h.ID.String(),
			UserID:     h.UserID.String(),
			SurahID:    h.SurahID,
			AyahNumber: h.AyahNumber,
			PageNumber: h.PageNumber,
			Mode:       h.Mode,
			LastReadAt: h.LastReadAt,
		})
	}

	return responses, nil
}

func (s *syncService) BulkInsertActivities(userID string, req dto.BulkActivityRequest) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	var activities []model.ReadingActivity
	for _, item := range req.Activities {
		activities = append(activities, model.ReadingActivity{
			UserID:          uid,
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
