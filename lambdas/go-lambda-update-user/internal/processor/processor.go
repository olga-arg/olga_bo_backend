package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-update-user/internal/storage"
	"go-lambda-update-user/pkg/domain"
	"go-lambda-update-user/pkg/dto"
)

type Processor interface {
	UpdateUserCardLimits(ctx context.Context, newUser *domain.User) error
	GetUser(ctx context.Context, userID string) (*domain.User, error)
	ValidateUserInput(ctx context.Context, input *dto.UpdateLimitInput, request events.APIGatewayProxyRequest) (*domain.User, error)
}

type processor struct {
	storage *storage.UserRepository
}

func NewProcessor(storage *storage.UserRepository) Processor {
	return &processor{
		storage: storage,
	}
}

func (p *processor) UpdateUserCardLimits(ctx context.Context, newUser *domain.User) error {
	err := p.storage.UpdateUserCardLimit(newUser)
	if err != nil {
		return err
	}
	return nil
}

func (p *processor) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	user, err := p.storage.GetUserByID(userID)
	if err != nil {
		fmt.Println("Error getting user by ID", err.Error())
		return nil, err
	}
	return user, nil
}

func (p *processor) ValidateUserInput(ctx context.Context, input *dto.UpdateLimitInput, request events.APIGatewayProxyRequest) (*domain.User, error) {
	fmt.Println("Validating input")
	if err := json.Unmarshal([]byte(request.Body), &input); err != nil {
		return nil, fmt.Errorf("invalid request body: %s", err.Error())
	}
	if input.PurchaseLimit < 0 {
		return nil, fmt.Errorf("invalid purchase limit")
	}
	if input.MonthlyLimit < 0 {
		return nil, fmt.Errorf("invalid monthly limit")
	}
	user, err := p.GetUser(ctx, request.PathParameters["user_id"])
	if err != nil {
		fmt.Println("error getting user", err.Error())
		return nil, fmt.Errorf("failed to get user")
	}
	if input.PurchaseLimit > 0 && input.MonthlyLimit > 0 {
		if input.PurchaseLimit > input.MonthlyLimit {
			fmt.Println("purchase limit cannot be greater than monthly limit")
			return nil, fmt.Errorf("purchase limit cannot be greater than monthly limit")
		}
		user.PurchaseLimit = input.PurchaseLimit
		user.MonthlyLimit = input.MonthlyLimit
	} else if input.PurchaseLimit > 0 {
		actualMonthlyLimit := user.MonthlyLimit
		if input.PurchaseLimit > actualMonthlyLimit {
			fmt.Println("purchase limit cannot be greater than monthly limit")
			return nil, fmt.Errorf("purchase limit cannot be greater than monthly limit")
		}
		user.PurchaseLimit = input.PurchaseLimit
	} else if input.MonthlyLimit > 0 {
		actualPurchaseLimit := user.PurchaseLimit
		if input.MonthlyLimit < actualPurchaseLimit {
			fmt.Println("monthly limit cannot be less than purchase limit")
			return nil, fmt.Errorf("monthly limit cannot be less than purchase limit")
		}
		user.MonthlyLimit = input.MonthlyLimit
	}
	fmt.Println("Input validated successfully")
	return user, nil
}
