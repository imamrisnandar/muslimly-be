package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSettings struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    *uuid.UUID `gorm:"type:uuid;index:idx_user_settings_unique,unique"`
	DeviceID  *uuid.UUID `gorm:"type:uuid;index:idx_user_settings_unique,unique"`
	Key       string     `gorm:"index:idx_user_settings_unique,unique;not null"`
	Value     string     `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserSettings) TableName() string {
	return "user_settings"
}
