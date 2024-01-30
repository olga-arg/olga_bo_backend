package processor

import (
	"context"
	"go-lambda-me/internal/storage"
	"go-lambda-me/pkg/dto"
	"time"
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
	// Transformar las fechas a UTC-3 (hora de Argentina)
	for i := range payments {
		if !payments[i].CreatedDate.IsZero() {
			payments[i].CreatedDate = payments[i].CreatedDate.In(time.FixedZone("UTC-3", -3*60*60))
		}
	}

	return dto.NewOutput(userInformation, payments), nil
}
