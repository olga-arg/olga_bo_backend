package processor

import (
	"commons/domain"
	"commons/utils/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-update-user/pkg/dto"
)

type Processor interface {
	UpdateUser(ctx context.Context, newUser *domain.User, companyId string) error
	GetUser(ctx context.Context, userID, companyId string) (*domain.User, error)
	ValidateUserInput(ctx context.Context, input *dto.UpdateUserInput, request events.APIGatewayProxyRequest, companyId string) (*domain.User, error)
	ValidateUser(ctx context.Context, email, companyId string, allowedRoles []domain.UserRoles) (bool, error)
}

type processor struct {
	userStorage *db.UserRepository
}

func NewProcessor(storage *db.UserRepository) Processor {
	return &processor{
		userStorage: storage,
	}
}

func (p *processor) UpdateUser(ctx context.Context, newUser *domain.User, companyId string) error {
	err := p.userStorage.UpdateUser(newUser, companyId)
	if err != nil {
		return err
	}
	return nil
}

func (p *processor) GetUser(ctx context.Context, userID, companyId string) (*domain.User, error) {
	user, err := p.userStorage.GetUserByID(userID, companyId)
	if err != nil {
		fmt.Println("Error getting user by ID", err.Error())
		return nil, err
	}
	return user, nil
}

func (p *processor) ValidateUserInput(ctx context.Context, input *dto.UpdateUserInput, request events.APIGatewayProxyRequest, companyId string) (*domain.User, error) {
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

	fmt.Println("Getting user", request.PathParameters["user_id"])

	user, err := p.GetUser(ctx, request.PathParameters["user_id"], companyId)
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

	if input.Role != "" {
		role, err := domain.ParseUserRole(input.Role)
		if err != nil {
			return nil, fmt.Errorf("invalid role")
		}
		user.Role = role
	}

	// Check if the role number is inside the UserRoles enum
	if user.Role < 0 || user.Role > 3 {
		return nil, fmt.Errorf("invalid role")
	}

	fmt.Println("Input validated successfully")
	return user, nil
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
