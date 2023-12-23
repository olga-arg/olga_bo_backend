package dto

type UpdateUserInput struct {
	PurchaseLimit int   `json:"purchase_limit"`
	MonthlyLimit  int   `json:"monthly_limit"`
	IsAdmin       *bool `json:"is_admin,omitempty"`
}
