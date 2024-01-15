package dto

import (
	"go-lambda-update-payment/pkg/domain"
	"time"
)

type ConfirmationStatus int

const (
	Pending ConfirmationStatus = iota
	Confirmed
	Deleted
)

type Output struct {
	ID              string             `json:"id"`
	Amount          float32            `json:"amount"`
	ShopName        string             `json:"shop_name"`
	Cuit            string             `json:"cuit"`
	Date            string             `json:"date"`
	Time            string             `json:"time"`
	UserID          string             `json:"user_id"`
	Category        string             `json:"category"`
	receiptNumber   string             `json:"receiptNumber"`
	receiptType     string             `json:"receiptType"`
	Status          ConfirmationStatus `json:"status" default:"Pending"`
	ReceiptImageKey string             `json:"receiptImageKey"`
	CreatedDate     time.Time          `json:"created"`
}

func NewOutput(payment *domain.Payment) *Output {
	dtoPayment := Output{
		Amount:          payment.Amount,
		ShopName:        payment.ShopName,
		Cuit:            payment.Cuit,
		Date:            payment.Date,
		Time:            payment.Time,
		UserID:          payment.UserID,
		Category:        payment.Category,
		Status:          ConfirmationStatus(payment.Status),
		ReceiptImageKey: payment.ReceiptImageKey,
		CreatedDate:     payment.CreatedDate,
	}
	return &dtoPayment
}
