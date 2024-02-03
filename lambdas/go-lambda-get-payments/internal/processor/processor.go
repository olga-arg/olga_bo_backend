package processor

import (
	"commons/domain"
	"commons/utils/db"
	"context"
	"go-lambda-get-payments/pkg/dto"
)

type Processor interface {
	GetAllPayments(ctx context.Context, filter map[string]string, companyId string) (*dto.Output, error)
	ValidateUser(ctx context.Context, email, companyId string, allowedRoles []domain.UserRoles) (bool, error)
}

type processor struct {
	paymentStorage db.PaymentRepository
	userStorage    db.UserRepository
}

func NewProcessor(paymentStorage db.PaymentRepository, userStorage db.UserRepository) Processor {
	return &processor{
		paymentStorage: paymentStorage,
		userStorage:    userStorage,
	}
}

func (p *processor) GetAllPayments(ctx context.Context, filter map[string]string, companyId string) (*dto.Output, error) {
	payments, err := p.paymentStorage.GetAllPayments(filter, companyId)
	if err != nil {
		return nil, err
	}
	return dto.NewOutput(payments), nil
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
