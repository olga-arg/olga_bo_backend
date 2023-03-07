package dto

type CreateUserInput struct {
	Name    string `json:"name" validate:"required"`
	Surname string `json:"surname" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}
