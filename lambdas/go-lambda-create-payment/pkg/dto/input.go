package dto

import "time"

type CreatePaymentInput struct {
	ShopName        string    `json:"shopName"`
	Category        string    `json:"category"`
	Cuit            string    `json:"cuit"`
	Date            time.Time `json:"date"`
	Time            string    `json:"time"`
	ReceiptNumber   string    `json:"receiptNumber"`
	ReceiptType     string    `json:"receiptType"`
	Amount          float32   `json:"amount"`
	ReceiptImageKey string    `json:"receiptImageKey"`
}
