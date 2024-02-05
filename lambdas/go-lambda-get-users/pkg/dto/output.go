package dto

import (
	"commons/domain"
)

type Output struct {
	Users []domain.User `json:"users"`
}

func NewOutput(users []domain.User) *Output {
	var dtoUsers []domain.User
	for _, user := range users {
		var teams []*domain.Team // Change the type to a slice of pointers
		for _, team := range user.Teams {
			// Append a pointer to the new Team object
			teams = append(teams, &domain.Team{
				ID:   team.ID,
				Name: team.Name,
			})
		}
		dtoUsers = append(dtoUsers, domain.User{
			ID:              user.ID,
			Name:            user.Name,
			Surname:         user.Surname,
			Email:           user.Email,
			PurchaseLimit:   user.PurchaseLimit,
			MonthlyLimit:    user.MonthlyLimit,
			MonthlySpending: user.MonthlySpending,
			Status:          user.Status,
			Teams:           teams,
		})
	}
	return &Output{
		Users: dtoUsers,
	}
}
