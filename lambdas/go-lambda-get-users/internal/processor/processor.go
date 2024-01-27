package processor

import (
	"commons/utils/db"
	"context"
	"go-lambda-get-all-users/pkg/dto"
)

type Processor interface {
	GetAllUsers(ctx context.Context, filter map[string]string, companyId string) (*dto.Output, error)
}

type processor struct {
	userStorage *db.UserRepository
}

func NewProcessor(storage *db.UserRepository) Processor {
	return &processor{
		userStorage: storage,
	}
}

func (p *processor) GetAllUsers(ctx context.Context, filter map[string]string, companyId string) (*dto.Output, error) {
	users, err := p.userStorage.GetAllUsers(filter, companyId)
	if err != nil {
		return nil, err
	}
	return dto.NewOutput(users), nil
}
