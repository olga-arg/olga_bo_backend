package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/badoux/checkmail"
	"github.com/pkg/errors"
	"go-lambda-create-user/internal/storage"
	"go-lambda-create-user/pkg/domain"
	"go-lambda-create-user/pkg/dto"
	
)

type Processor interface {
	CreateUser(ctx context.Context, input *dto.CreateUserInput) error
	ValidateUserInput(ctx context.Context, input *dto.CreateUserInput, request events.APIGatewayProxyRequest) error
}

type processor struct {
	storage storage.UserRepository
}

func New(s storage.UserRepository) Processor {
	return &processor{
		storage: s,
	}
}

func (p *processor) CreateUser(ctx context.Context, input *dto.CreateUserInput) error {
	// Checks if the email already exists
	exists, err := p.storage.EmailAlreadyExists(input.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	// Creates a new user. New user takes a name and email and returns a user struct
	user, err := domain.NewUser(input.Name, input.Surname, input.Email)
	if err != nil {
		fmt.Println("Error creating user: ", err)
		return err
	}
	// Saves the user to the database if it doesn't already exist
	if err := p.storage.Save(user); err != nil {
		fmt.Println("Error saving user: ", err)
		return err
	}
	// Returns
	return nil
}

func (p *processor) ValidateUserInput(ctx context.Context, input *dto.CreateUserInput, request events.APIGatewayProxyRequest) error {
	fmt.Println("Validating input")
	if request.Body == "" || len(request.Body) < 1 {
		return fmt.Errorf("missing request body")
	}
	if err := json.Unmarshal([]byte(request.Body), &input); err != nil {
		return fmt.Errorf("invalid request body: %s", err.Error())
	}
	if input.Name == "" {
		return fmt.Errorf("name is required")
	}
	if input.Surname == "" {
		return fmt.Errorf("surname is required")
	}
	if input.Email == "" {
		return fmt.Errorf("email is required")
	}
	if request.Body == "" || len(request.Body) < 1 {
		return fmt.Errorf("missing request body")
	}
	if len(input.Name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}
	if len(input.Surname) < 2 {
		return fmt.Errorf("surname must be at least 2 characters")
	}
	if len(input.Name) > 50 {
		return fmt.Errorf("name must be less than 50 characters")
	}
	if len(input.Surname) > 50 {
		return fmt.Errorf("surname must be less than 50 characters")
	}
	err := checkmail.ValidateFormat(input.Email)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}
	return nil
}
