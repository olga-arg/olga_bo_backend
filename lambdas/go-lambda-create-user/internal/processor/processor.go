package processor

import (
	"context"
	"go-lambda-create-user/internal/storage"
	"go-lambda-create-user/pkg/domain"
	"go-lambda-create-user/pkg/dto"
)

type Processor interface {
	CreateUser(ctx context.Context, input *dto.UserInput) (*dto.UserOutput, error)
}

type processor struct {
	storage storage.Storage
}

func New(s storage.Storage) Processor {
	return &processor{
		storage: s,
	}
}

func (p *processor) CreateUser(ctx context.Context, input *dto.UserInput) (*dto.UserOutput, error) {
	user := domain.NewUser(input.Name, input.Email)
	if err := p.storage.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return &dto.UserOutput{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
