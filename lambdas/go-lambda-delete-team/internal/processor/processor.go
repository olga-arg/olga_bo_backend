package processor

import (
	"commons/utils/db"
	"context"
	"fmt"
)

type Processor interface {
	DeleteTeam(ctx context.Context, teamID, companyId string) error
}

type processor struct {
	teamStorage *db.TeamRepository
}

func NewProcessor(storage *db.TeamRepository) Processor {
	return &processor{
		teamStorage: storage,
	}
}

func (p *processor) DeleteTeam(ctx context.Context, teamID, companyId string) error {
	fmt.Println("Deleting team in storage")
	err := p.teamStorage.DeleteTeam(teamID, companyId)
	if err != nil {
		fmt.Println("error", err.Error())
		return err
	}
	fmt.Println("Team deleted proc")
	return nil
}
