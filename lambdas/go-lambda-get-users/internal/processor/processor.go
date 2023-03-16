package processor

import (
	"context"
	"go-lambda-get-all-users/internal/storage"
	"go-lambda-get-all-users/pkg/dto"
)

type Processor interface {
	GetAllUsers(ctx context.Context, filter map[string]string) (*dto.Output, error)
}

type processor struct {
	storage *storage.UserRepository
}

func NewProcessor(storage *storage.UserRepository) Processor {
	return &processor{
		storage: storage,
	}
}

func (p *processor) GetAllUsers(ctx context.Context, filter map[string]string) (*dto.Output, error) {
	users, err := p.storage.GetAllUsers(filter)
	if err != nil {
		return nil, err
	}
	return dto.NewOutput(users), nil
}
