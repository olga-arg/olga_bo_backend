package processor

import (
	"commons/utils/db"
	"context"
	"go-lambda-get-teams/pkg/dto"
)

type Processor interface {
	GetAllTeams(ctx context.Context, filter map[string]string, companyId string) (*dto.Output, error)
}

type processor struct {
	teamStorage *db.TeamRepository
}

func NewProcessor(storage *db.TeamRepository) Processor {
	return &processor{
		teamStorage: storage,
	}
}

func (p *processor) GetAllTeams(ctx context.Context, filter map[string]string, companyId string) (*dto.Output, error) {
	teams, err := p.teamStorage.GetAllTeams(filter, companyId)
	if err != nil {
		return nil, err
	}
	//teams, err = p.teamStorage.GetAllReviewers(teams, companyId)
	return dto.NewOutput(teams), nil
}
