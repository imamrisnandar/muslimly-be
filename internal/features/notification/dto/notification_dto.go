package dto

type RegisterDeviceRequest struct {
	FCMToken string `json:"fcm_token" validate:"required"`
	Platform string `json:"platform" validate:"oneof=android ios web"`

	// Optional Device Info
	DeviceModel     string `json:"device_model,omitempty"`
	DeviceOSVersion string `json:"device_os_version,omitempty"`
	AppVersion      string `json:"app_version,omitempty"`

	// Optional Location Info
	CountryCode string `json:"country_code,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
}
