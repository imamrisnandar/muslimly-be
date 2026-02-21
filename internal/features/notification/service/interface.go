package service

import (
	"muslimly-be/internal/features/notification/dto"
	"muslimly-be/internal/features/notification/model"
)

type NotificationService interface {
	RegisterDevice(userID string, req dto.RegisterDeviceRequest) (*model.UserDevice, error)
	SendDailyReminder() error
}
