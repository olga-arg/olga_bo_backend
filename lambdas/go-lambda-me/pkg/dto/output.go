package dto

import (
	"commons/domain"
)

type UserInformation struct {
	Email           string  `json:"email"`
	Name            string  `json:"name"`
	Surname         string  `json:"surname"`
	FullName        string  `json:"full_name"`
	PurchaseLimit   int     `json:"purchase_limit"`
	MonthlyLimit    int     `json:"monthly_limit"`
	MonthlySpending float32 `json:"monthly_spending"`
}

type Output struct {
	UserInformation UserInformation  `json:"user_information"`
	Payments        []domain.Payment `json:"payments"`
}

func NewOutput(userInformation domain.User, payments []domain.Payment) *Output {
	var dtoPayments []domain.Payment
	for _, payment := range payments {
		dtoPayments = append(dtoPayments, domain.Payment{
			Amount:          payment.Amount,
			ShopName:        payment.ShopName,
			Category:        payment.Category,
			Status:          payment.Status,
			ReceiptImageKey: payment.ReceiptImageKey,
			CreatedDate:     payment.CreatedDate,
		})
	}
	var dtoUserInformation = UserInformation{
		Email:           userInformation.Email,
		Name:            userInformation.Name,
		Surname:         userInformation.Surname,
		FullName:        userInformation.FullName,
		PurchaseLimit:   userInformation.PurchaseLimit,
		MonthlyLimit:    userInformation.MonthlyLimit,
		MonthlySpending: userInformation.MonthlySpending,
	}
	return &Output{
		UserInformation: dtoUserInformation,
		Payments:        dtoPayments,
	}
}
