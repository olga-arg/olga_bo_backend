package processor

import (
	"context"
	"go-lambda-create-team/internal/storage"
	"go-lambda-create-team/pkg/domain"
	"go-lambda-create-team/pkg/dto"
	"log"
)

type Processor interface {
	CreateTeam(ctx context.Context, input *dto.CreateTeamInput) error
}

type processor struct {
	storage storage.TeamRepository
}

func New(s storage.TeamRepository) Processor {
	return &processor{
		storage: s,
	}
}

func (p *processor) CreateTeam(ctx context.Context, input *dto.CreateTeamInput) error {
	// Creates a new team
	team, err := domain.NewTeam(input.TeamName, input.ReviewerId, input.AnnualBudget, &p.storage) //input.Employees)
	if err != nil {
		log.Println("Error creating team: ", err)
		return err
	}
	// Saves the team to the database if it doesn't already exist
	if err := p.storage.Save(team); err != nil {
		log.Println("Error saving team: ", err)
		return err
	}
	// Returns
	return nil
}
