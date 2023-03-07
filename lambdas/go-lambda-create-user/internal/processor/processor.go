package processor

import (
	"context"
	"github.com/pkg/errors"
	"go-lambda-create-user/internal/storage"
	"go-lambda-create-user/pkg/domain"
	"go-lambda-create-user/pkg/dto"
)

type Processor interface {
	CreateUser(ctx context.Context, input *dto.CreateUserInput) (*dto.CreateUserOutput, error)
}

type processor struct {
	storage storage.UserRepository
}

func New(s storage.UserRepository) Processor {
	return &processor{
		storage: s,
	}
}

func (p *processor) CreateUser(ctx context.Context, input *dto.CreateUserInput) (*dto.CreateUserOutput, error) {
	exists, err := p.storage.EmailAlreadyExists(input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}
	user, _ := domain.NewUser(input.Name, input.Email)
	if err := p.storage.Save(user); err != nil {
		return nil, err
	}
	return &dto.CreateUserOutput{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
