package domain

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Team struct {
	ID              string             `json:"id"`
	Name            string             `json:"name"`
	Users           Users              `json:"users"`
	Reviewer        User               `json:"reviewer"`
	ReviewerId      string             `json:"reviewer_id"`
	MonthlySpending int                `json:"monthly_spending"`
	AnnualBudget    int                `json:"annual_budget"`
	Status          ConfirmationStatus `json:"status" default:"Pending"`
	CreatedDate     time.Time          `json:"created_date"`
}

type Teams []Team

type DbTeam struct {
	ID                  string             `json:"id"`
	Name                string             `json:"name"`
	ReviewerId          string             `json:"reviewer_id"`
	MonthlySpending     int                `json:"monthly_spending"`
	AnnualBudget        int                `json:"annual_budget"`
	Status              ConfirmationStatus `json:"status" default:"Pending"`
	CreatedDate         time.Time          `json:"created_date"`
	UserId              string             `json:"users.id"`
	UserName            string             `json:"users.name";gorm:"column:users.name"`
	UserSurname         string             `json:"users.surname"`
	UserEmail           string             `json:"users.email";gorm:"column:users.name"`
	UserMonthlySpending float32            `json:"users.monthly_spending"`
}

type DbTeams []DbTeam

func NewTeam(name, reviewer string, budget int) (*Team, error) {
	var team Team
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("error generating uuid: ", err)
		return nil, err
	}
	team.ID = id.String()
	team.Name = name
	team.ReviewerId = reviewer
	team.AnnualBudget = budget
	team.Status = Created
	team.CreatedDate = time.Now()
	return &team, nil
}
