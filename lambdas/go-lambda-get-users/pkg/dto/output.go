package dto

import (
	"go-lambda-get-all-users/pkg/domain"
)

type Team struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID              string                    `json:"id"`
	Name            string                    `json:"name"`
	Surname         string                    `json:"surname"`
	Email           string                    `json:"email"`
	PurchaseLimit   int                       `json:"purchase_limit" default:"0"`
	MonthlyLimit    int                       `json:"monthly_limit" default:"0"`
	MonthlySpending float32                   `json:"monthly_spending" default:"0"`
	IsAdmin         bool                      `json:"isAdmin" default:"false"`
	Status          domain.ConfirmationStatus `json:"status" default:"Pending"`
	Teams           []Team                    `json:"teams"`
}

type Output struct {
	Users []User `json:"users"`
}

// From domain.Users ([]User) to dto.Output (Output)
func NewOutput(users []domain.User) *Output {
	var dtoUsers []User
	for _, user := range users {
		var teams []Team
		for _, team := range user.Teams {
			teams = append(teams, Team{
				ID:   team.ID,
				Name: team.Name,
			})
		}
		dtoUsers = append(dtoUsers, User{
			ID:              user.ID,
			Name:            user.Name,
			Surname:         user.Surname,
			Email:           user.Email,
			PurchaseLimit:   user.PurchaseLimit,
			MonthlyLimit:    user.MonthlyLimit,
			MonthlySpending: user.MonthlySpending,
			IsAdmin:         user.IsAdmin,
			Status:          user.Status,
			Teams:           teams,
		})
	}
	return &Output{
		Users: dtoUsers,
	}
}
