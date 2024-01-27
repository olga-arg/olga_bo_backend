package dto

import (
	"commons/domain"
)

type Output struct {
	Team domain.Team `json:"team"`
}

// From domain.Teams ([]Team) to dto.Output (Output)
func NewOutput(team *domain.Team) *Output {
	return &Output{
		Team: domain.Team{
			ID:           team.ID,
			Name:         team.Name,
			ReviewerId:   team.ReviewerId,
			AnnualBudget: team.AnnualBudget,
			Status:       team.Status,
		},
	}
}
