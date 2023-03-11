package processor

import (
	"context"
	"go-lambda-create-user/internal/storage"
	"go-lambda-create-user/pkg/domain"
	"go-lambda-create-user/pkg/dto"
)

type Processor interface {
	CreateUser(ctx context.Context, input *dto.CreateUserInput) (*dto.CreateUserOutput, error)
}

type processor struct {
	storage storage.UserRepository
}

func New(s storage.UserRepository) Processor {
	return &processor{
		storage: s,
	}
}

func (p *processor) CreateUser(ctx context.Context, input *dto.CreateUserInput) (*dto.CreateUserOutput, error) {
	// Checks if the email already exists
	//exists, err := p.storage.EmailAlreadyExists(input.Email)
	//if err != nil {
	//	return nil, err
	//}
	//if exists {
	//	return nil, errors.New("email already exists")
	//}

	// Creates a new user. New user takes a name and email and returns a user struct
	user, _ := domain.NewUser(input.Name, input.Surname, input.Email)
	// Saves the user to the database if it doesn't already exist
	if err := p.storage.Save(user); err != nil {
		return nil, err
	}
	// Returns the user
	return &dto.CreateUserOutput{
		ID:      user.ID,
		Name:    user.Name,
		Surname: user.Surname,
		Email:   user.Email,
		Limit:   user.Limit,
		IsAdmin: user.IsAdmin,
		Status:  user.Status,
	}, nil
}
