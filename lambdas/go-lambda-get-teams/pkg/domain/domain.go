package domain

import "time"

type ConfirmationStatus int

const (
	Created ConfirmationStatus = iota
	Deleted
	Pending
)

type Team struct {
	ID           string             `json:"id"`
	CompanyID    string             `json:"company"`
	TeamName     string             `json:"name"`
	Users        []User             `gorm:"many2many:user_teams;"`
	ReviewerId   string             `json:"reviewer_id"`
	AnnualBudget int                `json:"annual_budget"`
	Status       ConfirmationStatus `json:"status" default:"Pending"`
	CreatedDate  time.Time          `json:"created_date"`
}

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
	IsAdmin         bool               `json:"isAdmin" default:"false"`
	Status          ConfirmationStatus `json:"status" default:"Pending"`
	CreatedDate     time.Time          `json:"created_date"`
	Teams           []Team             `gorm:"many2many:user_teams;"`
}

type Teams []Team

func ParseConfirmationStatus(s string) ConfirmationStatus {
	switch s {
	case "Pending":
		return Pending
	case "Created":
		return Created
	case "Deleted":
		return Deleted
	default:
		return Created
	}
}