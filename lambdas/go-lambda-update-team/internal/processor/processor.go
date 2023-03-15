package processor

import (
	"context"
	"fmt"
	"go-lambda-update-team/internal/storage"
	"go-lambda-update-team/pkg/dto"
)

type Processor interface {
	UpdateTeam(ctx context.Context, teamID string, newTeam *dto.UpdateTeamRequest) error
}

type processor struct {
	storage *storage.TeamRepository
}

func NewProcessor(storage *storage.TeamRepository) Processor {
	return &processor{
		storage: storage,
	}
}

func (p *processor) UpdateTeam(ctx context.Context, teamID string, newTeam *dto.UpdateTeamRequest) error {
	if newTeam.AnnualBudget < 0 {
		return fmt.Errorf("annual budget must be greater than 0")
	}
	fmt.Println("Updating team in storage")
	err := p.storage.UpdateTeamBudget(teamID, newTeam)
	if err != nil {
		fmt.Println("error", err.Error())
		return err
	}
	fmt.Println("Team updated proc")
	return nil
}
