package repository

import (
	"muslimly-be/internal/features/app_config/model"

	"gorm.io/gorm"
)

type AppConfigRepository interface {
	GetHijriAdjustments() ([]model.HijriAdjustment, error)
	UpsertHijriAdjustment(month, adjustment int) error
}

type appConfigRepository struct {
	db *gorm.DB
}

func NewAppConfigRepository(db *gorm.DB) AppConfigRepository {
	return &appConfigRepository{
		db: db,
	}
}

func (r *appConfigRepository) GetHijriAdjustments() ([]model.HijriAdjustment, error) {
	var adjustments []model.HijriAdjustment
	// Fetch all adjustments
	err := r.db.Find(&adjustments).Error
	return adjustments, err
}

func (r *appConfigRepository) UpsertHijriAdjustment(month, adjustment int) error {
	var adj model.HijriAdjustment
	err := r.db.Where("hijri_month = ?", month).First(&adj).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new
			newAdj := model.HijriAdjustment{
				HijriMonth: month,
				Adjustment: adjustment,
			}
			return r.db.Create(&newAdj).Error
		}
		return err
	}

	// Update existing
	adj.Adjustment = adjustment
	return r.db.Save(&adj).Error
}
