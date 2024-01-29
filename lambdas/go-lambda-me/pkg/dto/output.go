package dto

import (
	"go-lambda-me/pkg/domain"
	"time"
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

type Payment struct {
	ID          string                    `json:"id"`
	Amount      float32                   `json:"amount"`
	ShopName    string                    `json:"shop_name"`
	CardId      string                    `json:"card_id"`
	Type        domain.PaymentType        `json:"payment_type"`
	UserID      string                    `json:"user_id"`
	Category    string                    `json:"category"`
	Label       string                    `json:"label"`
	Status      domain.ConfirmationStatus `json:"status" default:"Pending"`
	Receipt     string                    `json:"receipt"`
	CreatedDate time.Time                 `json:"created"`
}

type Output struct {
	UserInformation UserInformation `json:"user_information"`
	Payments        []Payment       `json:"payments"`
}

func NewOutput(userInformation domain.User, payments []domain.Payment) *Output {
	var dtoPayments []Payment
	for _, payment := range payments {
		dtoPayments = append(dtoPayments, Payment{
			Amount:      payment.Amount,
			ShopName:    payment.ShopName,
			Type:        payment.Type,
			Category:    payment.Category,
			Label:       payment.Label,
			Status:      payment.Status,
			Receipt:     payment.Receipt,
			CreatedDate: payment.CreatedDate,
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
