package dto

import (
	"commons/domain"
)

type Output struct {
	Payments domain.Payments `json:"payments"`
}

func NewOutput(payments domain.Payments) *Output {
	var dtoPayments domain.Payments
	for _, payment := range payments {
		dtoPayments = append(dtoPayments, domain.Payment{
			ID:              payment.ID,
			Amount:          payment.Amount,
			ShopName:        payment.ShopName,
			Cuit:            payment.Cuit,
			Date:            payment.Date,
			Time:            payment.Time,
			UserID:          payment.UserID,
			Category:        payment.Category,
			ReceiptNumber:   payment.ReceiptNumber,
			ReceiptType:     payment.ReceiptType,
			Status:          payment.Status,
			ReceiptImageKey: payment.ReceiptImageKey,
			CreatedDate:     payment.CreatedDate,
			User: domain.User{
				Name:            payment.User.Name,
				Surname:         payment.User.Surname,
				FullName:        payment.User.FullName,
				Email:           payment.User.Email,
				PurchaseLimit:   payment.User.PurchaseLimit,
				MonthlyLimit:    payment.User.MonthlyLimit,
				MonthlySpending: payment.User.MonthlySpending,
				IsAdmin:         payment.User.IsAdmin,
			},
		})
	}
	return &Output{
		Payments: dtoPayments,
	}
}
