package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReadingHistory struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"type:uuid;index:idx_user_surah,unique;not null"`
	SurahID    int       `gorm:"index:idx_user_surah,unique;not null"`
	AyahNumber int       `gorm:"not null"`
	PageNumber int       `gorm:"default:0"`
	Mode       string    `gorm:"default:'mushaf'"`     // 'mushaf' or 'list'
	LastReadAt time.Time `gorm:"autoCreateTime:false"` // Manually set or updated
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// Ensure unique history per surah for a user
func (ReadingHistory) TableName() string {
	return "reading_histories"
}
