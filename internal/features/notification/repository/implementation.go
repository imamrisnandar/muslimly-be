package repository

import (
	"muslimly-be/internal/features/notification/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeviceRepository interface {
	UpsertDevice(device *model.UserDevice) error
	GetDevicesByUserID(userID string) ([]model.UserDevice, error)
	GetAllDevices() ([]model.UserDevice, error)
}

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{db}
}

func (r *deviceRepository) UpsertDevice(device *model.UserDevice) error {
	device.LastActiveAt = time.Now()
	// Upsert on FCMToken. If exists, update UserID & LastActive.
	// This handles case where user logs in on same device with different account,
	// or same account refreshes token.
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "fcm_token"}},
		DoUpdates: clause.AssignmentColumns([]string{"user_id", "platform", "last_active_at", "updated_at"}),
	}).Create(device).Error
}

func (r *deviceRepository) GetDevicesByUserID(userID string) ([]model.UserDevice, error) {
	var devices []model.UserDevice
	err := r.db.Where("user_id = ?", userID).Find(&devices).Error
	return devices, err
}

func (r *deviceRepository) GetAllDevices() ([]model.UserDevice, error) {
	var devices []model.UserDevice
	err := r.db.Find(&devices).Error
	return devices, err
}
