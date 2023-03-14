package dto

import (
	"go-lambda-delete-team/pkg/domain"
)

type Team struct {
	ID           string                    `json:"id"`
	TeamName     string                    `json:"team_name"`
	ReviewerId   string                    `json:"reviewer_id"`
	AnnualBudget int                       `json:"annual_budget"`
	Status       domain.ConfirmationStatus `json:"status" default:"Created"`
}

type Output struct {
	Team Team `json:"team"`
}

// From domain.Teams ([]Team) to dto.Output (Output)
func NewOutput(team *domain.Team) *Output {
	return &Output{
		Team: Team{
			ID:           team.ID,
			TeamName:     team.TeamName,
			ReviewerId:   team.ReviewerId,
			AnnualBudget: team.AnnualBudget,
			Status:       team.Status,
		},
	}
}
