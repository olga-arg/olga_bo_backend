package processor

import (
	"context"
	"go-lambda-me/internal/storage"
	"go-lambda-me/pkg/dto"
)

type Processor interface {
	GetUserInformation(ctx context.Context, email string) (*dto.Output, error)
}

type processor struct {
	storage *storage.PaymentRepository
}

func NewProcessor(storage *storage.PaymentRepository) Processor {
	return &processor{
		storage: storage,
	}
}

func (p *processor) GetUserInformation(ctx context.Context, email string) (*dto.Output, error) {
	userInformation, err := p.storage.GetUserInformation(email)
	if err != nil {
		return nil, err
	}
	payments, err := p.storage.GetAllPayments(userInformation.ID)
	if err != nil {
		return nil, err
	}

	return dto.NewOutput(userInformation, payments), nil
}
