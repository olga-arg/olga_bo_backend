package dto

type CreateUserOutput struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Surname   string   `json:"surname"`
	Email     string   `json:"email"`
	IsAdmin   bool     `json:"isAdmin"`
	Teams     []string `json:"team"`
	Confirmed bool     `json:"confirmed"`
}
