package dto

import (
	"commons/domain"
)

type Output struct {
	Categories domain.Categories `json:"categories"`
}

func NewOutput(categories domain.Categories) *Output {
	return &Output{
		Categories: categories,
	}
}
