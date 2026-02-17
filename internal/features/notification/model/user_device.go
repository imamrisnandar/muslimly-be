package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserDevice struct {
	ID       uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID   *uuid.UUID `gorm:"type:uuid;index"` // Nullable for Guest
	FCMToken string     `gorm:"type:text;unique;not null"`
	Platform string     `gorm:"type:varchar(20);default:'android'"` // android, ios, web

	// Device Info
	DeviceModel     string `gorm:"type:varchar(100)"`
	DeviceOSVersion string `gorm:"type:varchar(50)"`
	AppVersion      string `gorm:"type:varchar(20)"`

	// Location Info
	CountryCode string `gorm:"type:varchar(10)"`
	Timezone    string `gorm:"type:varchar(50)"`

	LastActiveAt time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (UserDevice) TableName() string {
	return "user_devices"
}
