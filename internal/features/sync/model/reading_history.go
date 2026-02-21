package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReadingHistory struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     *uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_rh_user_surah,where:user_id IS NOT NULL"`     // Nullable for guest
	DeviceID   *uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_rh_device_surah,where:device_id IS NOT NULL"` // Links to user_devices.id
	SurahID    int        `gorm:"not null;uniqueIndex:idx_rh_user_surah;uniqueIndex:idx_rh_device_surah"`
	AyahNumber int        `gorm:"not null"`
	PageNumber int        `gorm:"default:0"`
	Mode       string     `gorm:"default:'mushaf'"`     // 'mushaf' or 'list'
	LastReadAt time.Time  `gorm:"autoCreateTime:false"` // Manually set or updated
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// Ensure unique history per surah for a user
func (ReadingHistory) TableName() string {
	return "reading_histories"
}
