package service

import (
	"muslimly-be/internal/features/app_config/dto"
	"muslimly-be/internal/features/app_config/repository"
	"muslimly-be/pkg/config"
)

type AppConfigService interface {
	GetAppConfig() (dto.AppConfigResponse, error)
}

type appConfigService struct {
	cfg  *config.Config
	repo repository.AppConfigRepository
}

func NewAppConfigService(cfg *config.Config, repo repository.AppConfigRepository) AppConfigService {
	return &appConfigService{
		cfg:  cfg,
		repo: repo,
	}
}

func (s *appConfigService) GetAppConfig() (dto.AppConfigResponse, error) {
	adjustments, err := s.repo.GetHijriAdjustments()
	if err != nil {
		return dto.AppConfigResponse{}, err
	}

	var dtos []dto.HijriAdjustmentDTO
	for _, adj := range adjustments {
		dtos = append(dtos, dto.HijriAdjustmentDTO{
			Month:      adj.HijriMonth,
			Adjustment: adj.Adjustment,
		})
	}

	return dto.AppConfigResponse{
		HijriAdjustment:  0, // Keep for backward compatibility or global default
		HijriAdjustments: dtos,
	}, nil
}
