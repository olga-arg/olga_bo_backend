package dto

import (
	"go-lambda-update-user/pkg/domain"
)

type Output struct {
	Name            string                    `json:"name"`
	Surname         string                    `json:"surname"`
	Email           string                    `json:"email"`
	PurchaseLimit   int                       `json:"purchase_limit" default:"0"`
	MonthlyLimit    int                       `json:"monthly_limit" default:"0"`
	MonthlySpending float32                   `json:"monthly_spending" default:"0"`
	IsAdmin         bool                      `json:"isAdmin" default:"false"`
	Status          domain.ConfirmationStatus `json:"status" default:"Pending"`
}

func NewOutput(user *domain.User) *Output {
	dtoUser := Output{
		Name:            user.Name,
		Surname:         user.Surname,
		Email:           user.Email,
		PurchaseLimit:   user.PurchaseLimit,
		MonthlyLimit:    user.MonthlyLimit,
		MonthlySpending: user.MonthlySpending,
		IsAdmin:         user.IsAdmin,
		Status:          user.Status,
	}
	return &dtoUser
}
