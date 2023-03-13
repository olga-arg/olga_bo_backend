package domain

import "time"

type ConfirmationStatus int

const (
	Pending ConfirmationStatus = iota
	Confirmed
	Deleted
)

type User struct {
	ID              string             `json:"id"`
	CompanyID       string             `json:"company"`
	Name            string             `json:"name"`
	Surname         string             `json:"surname"`
	FullName        string             `json:"full_name"`
	Email           string             `json:"email"`
	PurchaseLimit   int                `json:"purchase_limit" default:"0"`
	MonthlyLimit    int                `json:"monthly_limit" default:"0"`
	MonthlySpending float32            `json:"monthly_spending" default:"0"`
	LimitResetDate  time.Time          `json:"reset_date"`
	IsAdmin         bool               `json:"isAdmin" default:"false"`
	Status          ConfirmationStatus `json:"status" default:"Pending"`
	CreatedDate     time.Time          `json:"created_date"`
}

type Users []User

func ParseConfirmationStatus(s string) ConfirmationStatus {
	switch s {
	case "Pending":
		return Pending
	case "Confirmed":
		return Confirmed
	case "Deleted":
		return Deleted
	default:
		return Pending
	}
}
