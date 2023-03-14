package domain

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go-lambda-create-team/internal/storage"
	"go-lambda-create-team/pkg/dto"
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
	input := &dto.CreateTeamInput{
		TeamName:     name,
		ReviewerId:   reviewer,
		AnnualBudget: budget,
	}

	err := validateInput(input, teamRepository)
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

func validateInput(input *dto.CreateTeamInput, teamRepository *storage.TeamRepository) error {
	// Validate Team Name is unique
	var existingTeam Team
	if err := teamRepository.Db.Where("team_name = ?", input.TeamName).First(&existingTeam).Error; err == nil {
		return errors.New("team name already exists")
	}
	if len(input.TeamName) > 50 {
		return errors.New("name must be less than 50 characters")
	}
	if input.TeamName == "" {
		return errors.New("name is required")
	}
	return nil
}
