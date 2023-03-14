package processor

import (
	"context"
	"fmt"
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
	// Validate that team isn't already deleted
	if newTeam.Status == 1 {
		log.Println("Team is already deleted")
		return fmt.Errorf("team is already deleted")
	}

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
	// if team not found return error
	if team == nil {
		log.Println("Team not found")
		return nil, fmt.Errorf("team not found")
	}
	return team, nil
}
