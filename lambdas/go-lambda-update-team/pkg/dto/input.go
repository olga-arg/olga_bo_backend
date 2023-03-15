package dto

type UpdateTeamRequest struct {
	Name         string   `json:"name"`
	AnnualBudget int      `json:"annual_budget"`
	Users        []string `json:"users"`
}
