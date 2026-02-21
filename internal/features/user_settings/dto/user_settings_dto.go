package dto

type SettingItem struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type UpsertSettingsRequest struct {
	DeviceID string        `json:"device_id,omitempty"`
	Settings []SettingItem `json:"settings" validate:"required,dive"`
}

type SettingResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
