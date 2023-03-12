package dto

import (
	"go-lambda-get-all-users/pkg/domain"
)

type User struct {
	Name         string                    `json:"name"`
	Surname      string                    `json:"surname"`
	Email        string                    `json:"email"`
	AccountLimit int                       `json:"limit"`
	Teams        []string                  `json:"teams"`
	IsAdmin      bool                      `json:"isAdmin"`
	Status       domain.ConfirmationStatus `json:"status"`
}

type Output struct {
	Users []User `json:"users"`
}

// From domain.Users ([]User) to dto.Output (Output)
func NewOutput(users []domain.User) *Output {
	var dtoUsers []User
	for _, user := range users {
		dtoUsers = append(dtoUsers, User{
			Name:         user.Name,
			Surname:      user.Surname,
			Email:        user.Email,
			AccountLimit: user.AccountLimit,
			Teams:        user.Teams,
			IsAdmin:      user.IsAdmin,
			Status:       user.Status,
		})
	}
	return &Output{
		Users: dtoUsers,
	}
}
