package dto

import (
	"commons/domain"
)

func NewOutput(payment *domain.Payment) *domain.Payment {
	dtoPayment := domain.Payment{
		Amount:          payment.Amount,
		ShopName:        payment.ShopName,
		Cuit:            payment.Cuit,
		Date:            payment.Date,
		Time:            payment.Time,
		UserID:          payment.UserID,
		Category:        payment.Category,
		Status:          payment.Status,
		ReceiptImageKey: payment.ReceiptImageKey,
		CreatedDate:     payment.CreatedDate,
	}
	return &dtoPayment
}
