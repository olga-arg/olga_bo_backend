package dto

type CreateTeamInput struct {
	Name         string `json:"name" validate:"required"`
	ReviewerId   string `json:"reviewer_id" validate:"required"`
	AnnualBudget int    `json:"annual_budget" validate:"required"`
}
