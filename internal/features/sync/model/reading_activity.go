package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReadingActivity struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          *uuid.UUID `gorm:"type:uuid;index"`                 // Nullable for guest
	DeviceID        *uuid.UUID `gorm:"type:uuid;index"`                 // Links to user_devices.id
	Date            string     `gorm:"type:varchar(10);index;not null"` // YYYY-MM-DD
	DurationSeconds int        `gorm:"not null;default:0"`
	PageNumber      int        `gorm:"default:0"`
	SurahNumber     int        `gorm:"default:0"`
	StartAyah       int        `gorm:"default:0"`
	EndAyah         int        `gorm:"default:0"`
	TotalAyahs      int        `gorm:"default:0"`
	Mode            string     `gorm:"default:'page'"`
	Timestamp       int64      `gorm:"not null"` // Unix timestamp from mobile
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (ReadingActivity) TableName() string {
	return "reading_activities"
}
