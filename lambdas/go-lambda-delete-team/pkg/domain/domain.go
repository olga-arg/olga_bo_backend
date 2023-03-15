package domain

import "time"

type ConfirmationStatus int

const (
	Created ConfirmationStatus = iota
	Deleted
	Pending
)

type Team struct {
	ID        string `json:"id"`
	CompanyID string `json:"company"`
	TeamName  string `json:"name"`
	//TODO: Create Employees to Team relationship
	//Employees    []string           `json:"employees"`
	ReviewerId      string             `json:"reviewer_id"`
	AnnualBudget    int                `json:"annual_budget"`
	MonthlySpending float32            `json:"monthly_spending" default:"0"`
	Status          ConfirmationStatus `json:"status" default:"Pending"`
	CreatedDate     time.Time          `json:"created_date"`
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
