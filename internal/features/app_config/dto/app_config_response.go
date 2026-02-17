package dto

type HijriAdjustmentDTO struct {
	Month      int `json:"hijri_month"`
	Adjustment int `json:"adjustment"`
}

type AppConfigResponse struct {
	HijriAdjustment  int                  `json:"hijri_adjustment"`
	HijriAdjustments []HijriAdjustmentDTO `json:"hijri_adjustments"`
}
