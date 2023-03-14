package processor

import (
	"context"
	"go-lambda-delete-team/internal/storage"
	"go-lambda-delete-team/pkg/domain"
	"log"
)

type Processor interface {
	DeleteTeam(ctx context.Context, newTeam *domain.Team) error
	GetTeam(ctx context.Context, teamID string) (*domain.Team, error)
}

type processor struct {
	storage *storage.TeamRepository
}

func NewProcessor(storage *storage.TeamRepository) Processor {
	return &processor{
		storage: storage,
	}
}

func (p *processor) DeleteTeam(ctx context.Context, newTeam *domain.Team) error {
	err := p.storage.DeleteTeam(newTeam)
	if err != nil {
		return err
	}
	return nil
}

func (p *processor) GetTeam(ctx context.Context, teamID string) (*domain.Team, error) {
	team, err := p.storage.GetTeamByID(teamID)
	if err != nil {
		log.Println("Error getting team by ID", err.Error())
		return nil, err
	}
	return team, nil
}
