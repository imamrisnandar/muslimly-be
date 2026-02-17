package dto

type RegisterDeviceRequest struct {
	FCMToken string `json:"fcm_token" validate:"required"`
	Platform string `json:"platform" validate:"oneof=android ios web"`
}
