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
	ID              string             `json:"id"`
	Amount          float32            `json:"amount"`
	ShopName        string             `json:"shop_name"`
	CardId          string             `json:"card_id"`
	Type            PaymentType        `json:"payment_type"`
	UserID          string             `json:"user_id"`
	Category        string             `json:"category"`
	Label           string             `json:"label"`
	Status          ConfirmationStatus `json:"status" default:"Pending"`
	ReceiptImageKey string             `json:"receipt_image_key"`
	CreatedDate     time.Time          `json:"created"`
	User            User               `gorm:"foreignKey:user_id"`
}

type Payments []Payment

type User struct {
	ID              string  `json:"id"`
	CompanyID       string  `json:"company"`
	Name            string  `json:"name"`
	Surname         string  `json:"surname"`
	FullName        string  `json:"full_name"`
	Email           string  `json:"email"`
	PurchaseLimit   int     `json:"purchase_limit" default:"0"`
	MonthlyLimit    int     `json:"monthly_limit" default:"0"`
	MonthlySpending float32 `json:"monthly_spending" default:"0"`
	IsAdmin         bool    `json:"isAdmin" default:"false"`
}

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
