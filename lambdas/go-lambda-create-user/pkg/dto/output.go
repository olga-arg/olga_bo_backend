package dto
import "go-lambda-create-user/pkg/domain"

type CreateUserOutput struct {
	ID        string                    `json:"id"`
	Name      string                    `json:"name"`
	Surname   string                    `json:"surname"`
	Email     string                    `json:"email"`
	Limit     int                       `json:"limit"`
	IsAdmin   bool                      `json:"isAdmin"`
	Teams     []string                  `json:"team"`
	Status    domain.ConfirmationStatus `json:"status"`
}
