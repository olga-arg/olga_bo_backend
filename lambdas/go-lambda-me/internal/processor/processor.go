package processor

import (
	"commons/utils/db"
	"context"
	"go-lambda-me/pkg/dto"
	"time"
)

type Processor interface {
	GetUserInformation(ctx context.Context, email, companyId string) (*dto.Output, error)
}

type processor struct {
	paymentStorage db.PaymentRepository
	userStorage    db.UserRepository
}

func NewProcessor(paymentRepo db.PaymentRepository, userRepo db.UserRepository) Processor {
	return &processor{
		paymentStorage: paymentRepo,
		userStorage:    userRepo,
	}
}

func (p *processor) GetUserInformation(ctx context.Context, email, companyId string) (*dto.Output, error) {
	userInformation, err := p.userStorage.GetUserInformation(email, companyId)
	if err != nil {
		return nil, err
	}
	payments, err := p.paymentStorage.GetAllPayments(nil, companyId)
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
