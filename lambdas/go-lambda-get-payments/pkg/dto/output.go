package dto

import (
	"go-lambda-get-payments/pkg/domain"
	"time"
)

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
	Payments []Payment `json:"payments"`
}

func NewOutput(payments []domain.Payment) *Output {
	var dtoPayments []Payment
	for _, payment := range payments {
		dtoPayments = append(dtoPayments, Payment{
			Amount:   payment.Amount,
			ShopName: payment.ShopName,
			Type:     payment.Type,
			Category: payment.Category,
			Label:    payment.Label,
			Status:   payment.Status,
			Receipt:  payment.Receipt,
		})
	}
	return &Output{
		Payments: dtoPayments,
	}
}
