package domain

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go-lambda-create-team/internal/storage"
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

func NewTeam(name, reviewer string, budget int, teamRepository *storage.TeamRepository) (*Team, error) {
	err := validateInput(name, teamRepository)
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

func validateInput(name string, teamRepository *storage.TeamRepository) error {
	// Validate Team Name is unique
	var existingTeam Team
	if err := teamRepository.Db.Where("team_name = ?", name).First(&existingTeam).Error; err == nil {
		return errors.New("team name already exists")
	}
	if len(name) > 50 {
		return errors.New("name must be less than 50 characters")
	}
	return nil
}
