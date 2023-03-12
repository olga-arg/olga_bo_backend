package processor

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
		strFilter[k] = fmt.Sprintf("%v", v)
	}

	// Perform type conversions for the filter values
	limit, err := strconv.Atoi(strFilter["limit"])
	if err != nil {
		return nil, fmt.Errorf("invalid limit value: %s", strFilter["limit"])
	}
	strFilter["limit"] = strconv.Itoa(limit)

	isAdmin, err := strconv.ParseBool(strFilter["isAdmin"])
	if err != nil {
		return nil, fmt.Errorf("invalid isAdmin value: %s", strFilter["isAdmin"])
	}
	strFilter["isAdmin"] = strconv.FormatBool(isAdmin)

	statusInt, err := strconv.Atoi(strFilter["status"])
	if err != nil {
		return nil, fmt.Errorf("invalid status value: %s", strFilter["status"])
	}
	strFilter["status"] = strconv.Itoa(statusInt)

	// Use the GetAllUsers method of the UserRepository to retrieve all users with pagination and filtering
	items, err := p.storage.GetAllUsers(strFilter, []string{"name", "surname", "email", "limit", "isAdmin", "team", "status"})
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
