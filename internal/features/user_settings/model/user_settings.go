package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSettings struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Key       string    `gorm:"primaryKey"`
	Value     string    `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserSettings) TableName() string {
	return "user_settings"
}
