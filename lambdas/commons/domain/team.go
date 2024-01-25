package domain

import "time"

type Team struct {
	ID           string             `json:"id"`
	CompanyID    string             `json:"company"`
	Name         string             `json:"name"`
	Users        []*User            `gorm:"many2many:user_teams;"`
	ReviewerId   string             `json:"reviewer_id"`
	AnnualBudget int                `json:"annual_budget"`
	Status       ConfirmationStatus `json:"status" default:"Pending"`
	CreatedDate  time.Time          `json:"created_date"`
}
