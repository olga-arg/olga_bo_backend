package dto

import "time"

type UpdatePaymentInput struct {
	ShopName string    `json:"shop_name,omitempty"`
	Amount   *float32  `json:"amount,omitempty"`
	Category string    `json:"category,omitempty"`
	Date     time.Time `json:"date,omitempty"`
	Currency string    `json:"currency,omitempty"`
	Status   *int      `json:"status,omitempty"`
}
