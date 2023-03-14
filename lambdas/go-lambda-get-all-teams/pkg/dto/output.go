package dto

import (
	"go-lambda-get-all-teams/pkg/domain"
)

type Team struct {
	ID           string                    `json:"id"`
	TeamName     string                    `json:"team_name"`
	ReviewerId   string                    `json:"reviewer_id"`
	AnnualBudget int                       `json:"annual_budget"`
	Status       domain.ConfirmationStatus `json:"status" default:"Created"`
}

type Output struct {
	Teams []Team `json:"teams"`
}

// From domain.Teams ([]Team) to dto.Output (Output)
func NewOutput(teams []domain.Team) *Output {
	var dtoTeams []Team
	for _, team := range teams {
		dtoTeams = append(dtoTeams, Team{
			ID:           team.ID,
			TeamName:     team.TeamName,
			ReviewerId:   team.ReviewerId,
			AnnualBudget: team.AnnualBudget,
			Status:       team.Status,
		})
	}
	return &Output{
		Teams: dtoTeams,
	}
}
