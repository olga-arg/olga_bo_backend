package processor

import (
	"commons/utils/db"
	"context"
	"go-lambda-get-payments/pkg/dto"
)

type Processor interface {
	GetAllPayments(ctx context.Context, filter map[string]string, companyId string) (*dto.Output, error)
}

type processor struct {
	paymentStorage *db.PaymentRepository
}

func NewProcessor(storage *db.PaymentRepository) Processor {
	return &processor{
		paymentStorage: storage,
	}
}

func (p *processor) GetAllPayments(ctx context.Context, filter map[string]string, companyId string) (*dto.Output, error) {
	payments, err := p.paymentStorage.GetAllPayments(filter, companyId)
	if err != nil {
		return nil, err
	}
	return dto.NewOutput(payments), nil
}
