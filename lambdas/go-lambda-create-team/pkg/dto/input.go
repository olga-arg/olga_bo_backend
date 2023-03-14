package dto

type CreateTeamInput struct {
	TeamName     string   `json:"name" validate:"required"`
	Employees    []string `json:"employees" validate:"required"`
	ReviewerId   string   `json:"reviewer_id" validate:"required"`
	AnnualBudget int      `json:"annual_budget" validate:"required"`
}
