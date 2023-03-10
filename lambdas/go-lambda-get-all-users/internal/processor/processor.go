package processor

import (
	"context"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"go-lambda-get-all-users/internal/storage"
	"go-lambda-get-all-users/pkg/dto"
)

type Processor interface {
	GetAllUsers(ctx context.Context, filter string) (*dto.Output, error)
}

type processor struct {
	storage *storage.UserRepository
}

func NewProcessor(storage *storage.UserRepository) Processor {
	return &processor{
		storage: storage,
	}
}

func (p *processor) GetAllUsers(ctx context.Context, filter string) (*dto.Output, error) {
	// Use the GetAllUsers method of the UserRepository to retrieve all users with pagination and filtering
	items, err := p.storage.GetAllUsers(filter, []string{"name", "surname", "email", "limit", "isAdmin", "team", "status"})
	if err != nil {
		return nil, err
	}

	// Convert each item to a *dynamodb.AttributeValue object
	attrValues := make([]*dynamodb.AttributeValue, 0, len(items))
	for _, item := range items {
		av := &dynamodb.AttributeValue{}
		av.M = item
		attrValues = append(attrValues, av)
	}

	// Convert each attribute value to a User object using the UnmarshalUsers function from the dto package
	users, err := dto.UnmarshalUsers(attrValues)
	if err != nil {
		return nil, err
	}

	return &dto.Output{
		Users: users,
	}, nil
}
