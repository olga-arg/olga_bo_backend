package dto

import (
	"commons/domain"
)

type Output struct {
	Teams domain.Teams `json:"teams"`
}

// From domain.Teams ([]Team) to dto.Output (Output)
func NewOutput(teams []domain.Team) *Output {
	var dtoTeams domain.Teams
	for _, team := range teams {
		var users domain.Users
		for _, user := range team.Users {
			users = append(users, domain.User{
				ID:      user.ID,
				Name:    user.Name,
				Surname: user.Surname,
				Email:   user.Email,
			})
		}
		dtoTeams = append(dtoTeams, domain.Team{
			ID:   team.ID,
			Name: team.Name,
			Reviewer: domain.User{
				ID:      team.Reviewer.ID,
				Name:    team.Reviewer.Name,
				Surname: team.Reviewer.Surname,
				Email:   team.Reviewer.Email,
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
