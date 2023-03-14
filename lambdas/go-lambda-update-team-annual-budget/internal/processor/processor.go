package processor

import (
	"context"
	"fmt"
	"go-lambda-update-team-annual-budget/internal/storage"
)

type Processor interface {
	UpdateTeamBudget(ctx context.Context, teamID string, annualBudget int) error
}

type processor struct {
	storage *storage.TeamRepository
}

func NewProcessor(storage *storage.TeamRepository) Processor {
	return &processor{
		storage: storage,
	}
}

func (p *processor) UpdateTeamBudget(ctx context.Context, teamID string, annualBudget int) error {
	fmt.Println("Updating team in storage")
	err := p.storage.UpdateTeamBudget(teamID, annualBudget)
	if err != nil {
		fmt.Println("error", err.Error())
		return err
	}
	fmt.Println("Team updated proc")
	return nil
}
