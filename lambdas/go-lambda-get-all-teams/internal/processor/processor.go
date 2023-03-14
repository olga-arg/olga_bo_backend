package processor

import (
	"context"
	"go-lambda-get-all-teams/internal/storage"
	"go-lambda-get-all-teams/pkg/dto"
)

type Processor interface {
	GetAllTeams(ctx context.Context, filter map[string]string) (*dto.Output, error)
}

type processor struct {
	storage *storage.TeamRepository
}

func NewProcessor(storage *storage.TeamRepository) Processor {
	return &processor{
		storage: storage,
	}
}

func (p *processor) GetAllTeams(ctx context.Context, filter map[string]string) (*dto.Output, error) {
	teams, err := p.storage.GetAllTeams(filter)
	if err != nil {
		return nil, err
	}
	return dto.NewOutput(teams), nil
}
