package processor

import (
	"commons/domain"
	"commons/utils/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-create-category/pkg/dto"
)

type Processor interface {
	CreateCategory(ctx context.Context, input *dto.CreateCategoryInput, companyId string) error
	ValidateCategoryInput(ctx context.Context, input *dto.CreateCategoryInput, request events.APIGatewayProxyRequest) error
}

type processor struct {
	categoryStorage db.CategoryStorage
}

func New(c db.CategoryStorage) Processor {
	return &processor{
		categoryStorage: c,
	}
}

func (p *processor) CreateCategory(ctx context.Context, input *dto.CreateCategoryInput, companyId string) error {
	// Creates a new category.
	category, err := domain.NewCategory(companyId, input.Name, input.Color, input.Icon)
	if err != nil {
		fmt.Println("Error creating category: ", err)
		return err
	}

	// Saves the company to the database if it doesn't already exist
	if err := p.categoryStorage.Save(category); err != nil {
		fmt.Println("Error saving category: ", err)
		return err
	}

	fmt.Println("Category created successfully")
	return nil
}

func (p *processor) ValidateCategoryInput(ctx context.Context, input *dto.CreateCategoryInput, request events.APIGatewayProxyRequest) error {
	fmt.Println("Validating input")
	if request.Body == "" || len(request.Body) < 1 {
		return fmt.Errorf("missing request body")
	}
	if err := json.Unmarshal([]byte(request.Body), &input); err != nil {
		return fmt.Errorf("invalid request body: %s", err.Error())
	}
	if input.Name == "" {
		return fmt.Errorf("name is required")
	}
	if input.Color == "" {
		return fmt.Errorf("color is required")
	}
	if input.Icon == "" {
		return fmt.Errorf("icon is required")
	}
	return nil
}
