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
	ID           string                    `json:"id"`
	TeamName     string                    `json:"team_name"`
	ReviewerId   string                    `json:"reviewer_id"`
	AnnualBudget int                       `json:"annual_budget"`
	Status       domain.ConfirmationStatus `json:"status" default:"Created"`
	Users        []User                    `json:"users"`
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
			ID:           team.ID,
			TeamName:     team.TeamName,
			ReviewerId:   team.ReviewerId,
			AnnualBudget: team.AnnualBudget,
			Status:       team.Status,
			Users:        users,
		})
	}
	return &Output{
		Teams: dtoTeams,
	}
}