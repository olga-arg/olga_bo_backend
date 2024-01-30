package processor

import (
	"commons/domain"
	"commons/utils/db"
	"context"
	"fmt"
)

type Processor interface {
	UpdateTeam(ctx context.Context, teamID string, newTeam *domain.UpdateTeamRequest, companyId string) error
}

type processor struct {
	teamStorage *db.TeamRepository
}

func NewProcessor(storage *db.TeamRepository) Processor {
	return &processor{
		teamStorage: storage,
	}
}

func (p *processor) UpdateTeam(ctx context.Context, teamID string, newTeam *domain.UpdateTeamRequest, companyId string) error {
	if newTeam.AnnualBudget < 0 {
		return fmt.Errorf("annual budget must be greater than 0")
	}
	fmt.Println("Updating team in storage")
	err := p.teamStorage.UpdateTeam(teamID, newTeam, companyId)
	if err != nil {
		fmt.Println("error", err.Error())
		return err
	}
	fmt.Println("Team updated proc")
	return nil
}
