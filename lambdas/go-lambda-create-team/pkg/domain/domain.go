package domain

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log"
	"time"
)

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
	// TODO: Create Employees to Team relationship
	//Employees    []string           `json:"employees"`
	ReviewerId   string             `json:"reviewer_id"`
	AnnualBudget int                `json:"annual_budget"`
	Status       ConfirmationStatus `json:"status" default:"Pending"`
	CreatedDate  time.Time          `json:"created_date"`
}

func NewTeam(name, reviewer string, budget int) (*Team, error) {
	err := validateInput(name)
	if err != nil {
		log.Println("error validating input: ", err)
		return nil, err
	}
	var team Team
	id, err := uuid.NewUUID()
	if err != nil {
		log.Println("error generating uuid: ", err)
		return nil, err
	}
	team.ID = id.String()
	team.TeamName = name
	//team.Employees = employees
	team.ReviewerId = reviewer
	team.AnnualBudget = budget
	team.Status = Created
	team.CreatedDate = time.Now()
	return &team, nil
}

func validateInput(name string) error {
	if len(name) > 50 {
		return errors.New("name must be less than 50 characters")
	}
	return nil
}
