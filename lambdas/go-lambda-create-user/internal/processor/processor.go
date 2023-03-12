package processor

import (
	"context"
	"github.com/pkg/errors"
	"go-lambda-create-user/internal/storage"
	"go-lambda-create-user/pkg/domain"
	"go-lambda-create-user/pkg/dto"
	"log"
)

type Processor interface {
	CreateUser(ctx context.Context, input *dto.CreateUserInput) error
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
		log.Println("Error creating user: ", err)
		return err
	}
	// Saves the user to the database if it doesn't already exist
	if err := p.storage.Save(user); err != nil {
		log.Println("Error saving user: ", err)
		return err
	}
	// Returns
	return nil
}
