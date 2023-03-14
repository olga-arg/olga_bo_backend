package processor

import (
	"context"
	"fmt"
	"go-lambda-delete-team/internal/storage"
)

type Processor interface {
	DeleteTeam(ctx context.Context, teamID string) error
}

type processor struct {
	storage *storage.TeamRepository
}

func NewProcessor(storage *storage.TeamRepository) Processor {
	return &processor{
		storage: storage,
	}
}

func (p *processor) DeleteTeam(ctx context.Context, teamID string) error {
	fmt.Println("Deleting team in storage")
	err := p.storage.DeleteTeam(teamID)
	if err != nil {
		fmt.Println("error", err.Error())
		return err
	}
	fmt.Println("Team deleted proc")
	return nil
}
