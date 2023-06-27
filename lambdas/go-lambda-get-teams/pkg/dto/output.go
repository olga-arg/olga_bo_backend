package dto

import (
	"go-lambda-get-teams/pkg/domain"
)

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type Team struct {
	ID              string                    `json:"id"`
	Name            string                    `json:"name"`
	Reviewer        User                      `json:"reviewer"`
	AnnualBudget    int                       `json:"annual_budget"`
	MonthlySpending float32                   `json:"monthly_spending" default:"0"`
	Status          domain.ConfirmationStatus `json:"status" default:"Created"`
	Users           []User                    `json:"users"`
}

type Output struct {
	Teams []Team `json:"teams"`
}

// From domain.Teams ([]Team) to dto.Output (Output)
func NewOutput(teams []domain.Team) *Output {
	var dtoTeams []Team
	for _, team := range teams {
		var users []User
		for _, user := range team.Users {
			users = append(users, User{
				ID:      user.ID,
				Name:    user.Name,
				Surname: user.Surname,
			})
		}
		dtoTeams = append(dtoTeams, Team{
			ID:   team.ID,
			Name: team.Name,
			Reviewer: User{
				ID:      team.Reviewer.ID,
				Name:    team.Reviewer.Name,
				Surname: team.Reviewer.Surname,
			},
			AnnualBudget:    team.AnnualBudget,
			MonthlySpending: team.MonthlySpending,
			Status:          team.Status,
			Users:           users,
		})
	}
	return &Output{
		Teams: dtoTeams,
	}
}
