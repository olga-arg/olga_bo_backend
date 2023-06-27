package dto

type CreateTeamInput struct {
	TeamName     string `json:"team_name" validate:"required"`
	ReviewerId   string `json:"reviewer_id" validate:"required"`
	AnnualBudget int    `json:"annual_budget" validate:"required"`
}
