package processor

import (
	"context"
	"go-lambda-get-payments/internal/storage"
	"go-lambda-get-payments/pkg/dto"
)

type Processor interface {
	GetAllPayments(ctx context.Context, filter map[string]string) (*dto.Output, error)
}

type processor struct {
	storage *storage.PaymentRepository
}

func NewProcessor(storage *storage.PaymentRepository) Processor {
	return &processor{
		storage: storage,
	}
}

func (p *processor) GetAllPayments(ctx context.Context, filter map[string]string) (*dto.Output, error) {
	payments, err := p.storage.GetAllPayments(filter)
	if err != nil {
		return nil, err
	}
	return dto.NewOutput(payments), nil
}
