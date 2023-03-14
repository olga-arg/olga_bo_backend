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
	// Creates a new user. New user takes a name and email and returns a user struct
	user, err := domain.NewTeam(input.TeamName, input.ReviewerId, input.AnnualBudget, input.Employees)
	if err != nil {
		log.Println("Error creating user: ", err)
		return err
	}
	// Saves the user to the database if it doesn't already exist
	if err := p.storage.Save(user); err != nil {
		log.Println("Error saving user: ", err)
		return err
	}
	// Returns
	return nil
}
