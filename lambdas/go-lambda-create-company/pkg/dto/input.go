package dto

type CreateCompanyInput struct {
	CompanyName string `json:"company_name" validate:"required"`
	UserName    string `json:"name" validate:"required"`
	UserSurname string `json:"surname" validate:"required"`
	UserEmail   string `json:"email" validate:"required,email"`
}
