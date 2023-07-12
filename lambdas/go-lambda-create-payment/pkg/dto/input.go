package dto

type CreatePaymentInput struct {
	ShopName string  `json:"shopName"`
	Amount   float32 `json:"amount"`
	CardID   string  `json:"cardID"`
	UserID   string  `json:"userID"`
	Category string  `json:"category"`
	Receipt  string  `json:"receipt"`
}
