package domain

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

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
