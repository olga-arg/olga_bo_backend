package domain

import (
	"github.com/badoux/checkmail"
	"github.com/pkg/errors"
	"time"
)

type ConfirmationStatus int

const (
	Pending ConfirmationStatus = iota
	Confirmed
	Deleted
)

type User struct {
	ID              string             `json:"id"`
	CompanyID       string             `json:"company"`
	Name            string             `json:"name"`
	Surname         string             `json:"surname"`
	FullName        string             `json:"full_name"`
	Email           string             `json:"email"`
	PurchaseLimit   int                `json:"purchase_limit" default:"0"`
	MonthlyLimit    int                `json:"monthly_limit" default:"0"`
	MonthlySpending float32            `json:"monthly_spending" default:"0"`
	IsAdmin         bool               `json:"isAdmin" default:"false"`
	Status          ConfirmationStatus `json:"status" default:"Pending"`
	CreatedDate     time.Time          `json:"created_date"`
}

func NewUser(name, surname, email string) (*User, error) {
	err := validateInput(name, surname, email)
	if err != nil {
		return nil, err
	}
	var user User
	user.Name = name
	user.Surname = surname
	user.Email = email
	user.Status = Pending
	user.CreatedDate = time.Now()
	return &user, nil
}

func validateInput(name, surname, email string) error {
	if len(email) > 50 {
		return errors.New("email must be less than 50 characters")
	}
	if len(name) > 50 {
		return errors.New("name must be less than 50 characters")
	}
	if len(surname) > 50 {
		return errors.New("surname must be less than 50 characters")
	}
	err := checkmail.ValidateFormat(email)
	if err != nil {
		return errors.New("invalid email format")
	}
	return nil
}
