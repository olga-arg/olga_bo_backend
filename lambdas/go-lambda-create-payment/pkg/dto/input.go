package dto

type CreatePaymentInput struct {
	Label  string  `json:"label"`
	Amount float32 `json:"amount"`
}
