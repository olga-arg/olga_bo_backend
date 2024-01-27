package dto

import (
	"commons/domain"
)

func NewOutput(user *domain.User) *domain.User {
	dtoUser := domain.User{
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
