package processor

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go-lambda-get-all-users/internal/storage"
	"go-lambda-get-all-users/pkg/dto"
	"strconv"
)

type Processor interface {
	GetAllUsers(ctx context.Context, filter map[string]interface{}) (*dto.Output, error)
}

type processor struct {
	storage *storage.UserRepository
}

func NewProcessor(storage *storage.UserRepository) Processor {
	return &processor{
		storage: storage,
	}
}

func (p *processor) GetAllUsers(ctx context.Context, filter map[string]interface{}) (*dto.Output, error) {
	// Convert the input filter map to a filter map with string values
	strFilter := make(map[string]string)
	for k, v := range filter {
		switch v := v.(type) {
		case bool:
			strFilter[k] = strconv.FormatBool(v)
		case int:
			strFilter[k] = strconv.Itoa(v)
		default:
			strFilter[k] = fmt.Sprint(v)
		}
	}

	// Use the GetAllUsers method of the UserRepository to retrieve all users with pagination and filtering
	items, err := p.storage.GetAllUsers(strFilter, []string{"name", "surname", "email", "limit", "isAdmin", "team", "status"})
	if err != nil {
		return nil, err
	}

	// Convert each item to a User object using the UnmarshalUser function from the dto package
	var users []*dto.User
	for _, item := range items {
		user := &dto.User{}
		if err := dynamodbattribute.UnmarshalMap(item, user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &dto.Output{
		Users: users,
	}, nil
}
