package dto

type UpdateLimitInput struct {
	PurchaseLimit int `json:"purchase_limit"`
	MonthlyLimit  int `json:"monthly_limit"`
}
