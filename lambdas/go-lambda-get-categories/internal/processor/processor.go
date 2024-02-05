package processor

import (
	"commons/domain"
	"commons/utils/db"
	"context"
	"go-lambda-get-categories/pkg/dto"
)

type Processor interface {
	GetCategories(ctx context.Context, companyId string) (map[string]interface{}, error)
	ValidateUser(ctx context.Context, email, companyId string, allowedRoles []domain.UserRoles) (bool, error)
}

type processor struct {
	categoryStorage db.CategoryRepository
	userStorage     db.UserRepository
}

func NewProcessor(categoryStorage db.CategoryRepository, userStorage db.UserRepository) Processor {
	return &processor{
		categoryStorage: categoryStorage,
		userStorage:     userStorage,
	}
}

func (p *processor) GetCategories(ctx context.Context, companyId string) (map[string]interface{}, error) {
	categories, err := p.categoryStorage.GetCategories(companyId)
	if err != nil {
		return nil, err
	}
	return dto.NewOutput(categories), nil
}

func (p *processor) ValidateUser(ctx context.Context, email, companyId string, allowedRoles []domain.UserRoles) (bool, error) {
	// Validate user
	isAuthorized, err := p.userStorage.IsUserAuthorized(email, companyId, allowedRoles)
	if err != nil {
		return false, err
	}
	if isAuthorized {
		return true, nil
	}
	return false, nil
}
