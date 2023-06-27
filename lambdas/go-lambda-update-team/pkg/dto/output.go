package dto

import (
	"go-lambda-update-team/pkg/domain"
)

type Team struct {
	ID           string                    `json:"id"`
	Name         string                    `json:"name"`
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
			Name:         team.Name,
			ReviewerId:   team.ReviewerId,
			AnnualBudget: team.AnnualBudget,
			Status:       team.Status,
		},
	}
}
