package processor

import (
	"context"
	"go-lambda-update-card-limit/internal/storage"
	"go-lambda-update-card-limit/pkg/domain"
)

type Processor interface {
	UpdateUserCardLimits(ctx context.Context, userID string, purchaseLimit int, monthlyLimit int) (*domain.User, error)
	GetUser(ctx context.Context, userID string) (*domain.User, error)
}

type processor struct {
	storage *storage.UserRepository
}

func NewProcessor(storage *storage.UserRepository) Processor {
	return &processor{
		storage: storage,
	}
}

func (p *processor) UpdateUserCardLimits(ctx context.Context, userID string, purchaseLimit int, monthlyLimit int) (*domain.User, error) {
	user, err := p.storage.UpdateUserCardLimit(userID, purchaseLimit, monthlyLimit)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *processor) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	user, err := p.storage.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
