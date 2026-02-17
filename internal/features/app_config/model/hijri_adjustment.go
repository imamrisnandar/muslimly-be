package model

import (
	"time"

	"gorm.io/gorm"
)

type HijriAdjustment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	HijriMonth int `gorm:"not null;uniqueIndex" json:"hijri_month"` // 1-12
	Adjustment int `gorm:"not null;default:0" json:"adjustment"`    // Days offset (-2 to +2)
}

func (HijriAdjustment) TableName() string {
	return "hijri_adjustments"
}
