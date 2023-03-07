package domain

import (
	"github.com/badoux/checkmail"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type User struct {
	ID        string   `json:"id"`
	Name      string   `json:"name, required"`
	Surname   string   `json:"surname, required"`
	Email     string   `json:"email, required"`
	IsAdmin   bool     `json:"isAdmin" default:"false"`
	Teams     []string `json:"team" default:"[]"`
	Confirmed bool     `json:"confirmed" default:"false"`
}

func NewUser(name, surname, email string) (*User, error) {
	user, err := validateInput(name, surname, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func validateInput(name, surname, email string) (*User, error) {
	if len(email) > 50 {
		return nil, errors.New("email must be less than 50 characters")
	}
	if len(name) > 50 {
		return nil, errors.New("name must be less than 50 characters")
	}
	if len(surname) > 50 {
		return nil, errors.New("surname must be less than 50 characters")
	}
	err := checkmail.ValidateFormat(email)
	if err != nil {
		return nil, errors.New("invalid email format")
	}
	return &User{
		ID:      uuid.NewString(),
		Name:    name,
		Surname: surname,
		Email:   email,
	}, nil
}
