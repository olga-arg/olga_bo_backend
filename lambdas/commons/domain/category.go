package domain

import (
	"fmt"
	"github.com/google/uuid"
)

type Category struct {
	ID        string `json:"id"`
	CompanyId string `json:"companyId"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Icon      string `json:"icon"`
}

type Categories []Category

func NewCategory(companyId, name, color, icon string) (*Category, error) {
	var category Category
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("error generating uuid: ", err)
		return nil, err
	}
	category.ID = id.String()
	category.CompanyId = companyId
	category.Name = name
	category.Color = color
	category.Icon = icon
	return &category, nil
}
