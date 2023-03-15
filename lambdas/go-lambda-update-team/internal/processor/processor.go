package processor

import (
	"context"
	"fmt"
	"go-lambda-update-team/internal/storage"
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
	if annualBudget < 0 {
		return fmt.Errorf("annual budget must be greater than 0")
	}
	fmt.Println("Updating team in storage")
	err := p.storage.UpdateTeamBudget(teamID, annualBudget)
	if err != nil {
		fmt.Println("error", err.Error())
		return err
	}
	fmt.Println("Team updated proc")
	return nil
}
