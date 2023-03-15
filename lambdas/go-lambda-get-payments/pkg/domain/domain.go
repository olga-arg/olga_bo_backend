package domain

import (
	"time"
)

type ConfirmationStatus int

const (
	Pending ConfirmationStatus = iota
	Canceled
	Approved
)

type PaymentType int

const (
	Card = iota
	Cash
)

type Payment struct {
	ID          string             `json:"id"`
	Amount      float32            `json:"amount"`
	ShopName    string             `json:"shop_name"`
	CardId      string             `json:"card_id"`
	Type        PaymentType        `json:"payment_type"`
	UserID      string             `json:"user_id"`
	Category    string             `json:"category"`
	Status      ConfirmationStatus `json:"status" default:"Pending"`
	Receipt     string             `json:"receipt"`
	CreatedDate time.Time          `json:"created"`
}

type Payments []Payment

func ParseConfirmationStatus(s string) ConfirmationStatus {
	switch s {
	case "Pending":
		return Pending
	case "Canceled":
		return Canceled
	case "Approved":
		return Approved
	default:
		return Pending
	}
}

func ParsePaymentType(s string) PaymentType {
	switch s {
	case "Card":
		return Card
	case "Cash":
		return Cash
	default:
		return Card
	}
}
