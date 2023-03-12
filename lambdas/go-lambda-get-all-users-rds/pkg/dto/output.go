package dto

import (
	"go-lambda-get-all-users/pkg/domain"
)

type User struct {
	Name    string                    `json:"name"`
	Surname string                    `json:"surname"`
	Email   string                    `json:"email"`
	Limit   int                       `json:"limit"`
	IsAdmin bool                      `json:"isAdmin"`
	Status  domain.ConfirmationStatus `json:"status"`
}

type Output struct {
	Users []User `json:"users"`
}

// From domain.Users ([]User) to dto.Output (Output)
func NewOutput(users domain.Users) *Output {
	var dtoUsers []User
	for _, user := range users {
		dtoUsers = append(dtoUsers, User{
			Name:    user.Name,
			Surname: user.Surname,
			Email:   user.Email,
			Limit:   user.Limit,
			IsAdmin: user.IsAdmin,
			Status:  user.Status,
		})
	}
	return &Output{
		Users: dtoUsers,
	}
}
