package dto

type CreateCategoryInput struct {
	Name  string `json:"name" validate:"required"`
	Color string `json:"color" validate:"required"`
	Icon  string `json:"icon" validate:"required"`
}
