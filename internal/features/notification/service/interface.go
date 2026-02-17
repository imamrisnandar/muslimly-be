package service

import "muslimly-be/internal/features/notification/dto"

type NotificationService interface {
	RegisterDevice(userID string, req dto.RegisterDeviceRequest) error
	SendDailyReminder() error
}
