package dto

type UpdateTeamRequest struct {
	Name         string   `json:"name"`
	AnnualBudget int      `json:"annual_budget"`
	ReviewerId   string   `json:"reviewer_id"`
	AddUsers     []string `json:"add_users"`
	RemoveUsers  []string `json:"remove_users"`
}
