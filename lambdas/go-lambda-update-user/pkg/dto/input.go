package dto

import (
	"commons/domain"
)

type UpdateUserInput struct {
	PurchaseLimit int               `json:"purchase_limit"`
	MonthlyLimit  int               `json:"monthly_limit"`
	IsAdmin       *bool             `json:"is_admin,omitempty"`
	Role          *domain.UserRoles `json:"role"`
}
