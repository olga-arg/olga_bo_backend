package processor

import (
	"commons/utils/db"
	"context"
	"go-lambda-get-categories/pkg/dto"
)

type Processor interface {
	GetCategories(ctx context.Context, companyId string) (map[string]interface{}, error)
}

type processor struct {
	categoryStorage *db.CategoryStorage
}

func NewProcessor(storage *db.CategoryStorage) Processor {
	return &processor{
		categoryStorage: storage,
	}
}

func (p *processor) GetCategories(ctx context.Context, companyId string) (map[string]interface{}, error) {
	categories, err := p.categoryStorage.GetCategories(companyId)
	if err != nil {
		return nil, err
	}
	return dto.NewOutput(categories), nil
}
