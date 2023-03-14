package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-update-card-limit/internal/storage"
	"go-lambda-update-card-limit/pkg/domain"
	"go-lambda-update-card-limit/pkg/dto"
)

type Processor interface {
	UpdateUserCardLimits(ctx context.Context, userID string, purchaseLimit *int, monthlyLimit *int) (*domain.User, error)
	GetUser(ctx context.Context, userID string) (*domain.User, error)
	ValidateUserInput(ctx context.Context, input *dto.UpdateLimitInput, request events.APIGatewayProxyRequest) error
}

type processor struct {
	storage *storage.UserRepository
}

func NewProcessor(storage *storage.UserRepository) Processor {
	return &processor{
		storage: storage,
	}
}

func (p *processor) UpdateUserCardLimits(ctx context.Context, userID string, purchaseLimit *int, monthlyLimit *int) (*domain.User, error) {
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

func (p *processor) ValidateUserInput(ctx context.Context, input *dto.UpdateLimitInput, request events.APIGatewayProxyRequest) error {
	if err := json.Unmarshal([]byte(request.Body), &input); err != nil {
		return fmt.Errorf("invalid request body: %s", err.Error())
	}

	// Validate input
	if input.PurchaseLimit < 0 {
		return fmt.Errorf("invalid purchase limit")
	}
	if input.MonthlyLimit < 0 {
		return fmt.Errorf("invalid monthly limit")
	}

	if input.MonthlyLimit < input.PurchaseLimit {
		return fmt.Errorf("monthly limit cannot be less than purchase limit")
	}

	return nil
}
