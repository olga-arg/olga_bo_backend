package dto

import (
	"commons/domain"
)

func NewOutput(categories domain.Categories) map[string]interface{} {
	var response = make(map[string]interface{})
	for _, category := range categories {
		response[category.Name] = map[string]interface{}{
			"color": category.Color,
			"icon":  category.Icon,
		}
	}
	return response
}
