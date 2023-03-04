package application

import (
	"context"
	"go-lambda-create-user/internal/processor"
	"go-lambda-create-user/internal/storage"
)

type App struct {
	Processor processor.Processor
	storage   storage.Storage
}

func New(ctx context.Context) (*App, error) {
	s, err := storage.NewDynamoDB()
	if err != nil {
		return nil, err
	}
	p := processor.New(s)
	return &App{
		Processor: p,
	}, nil
}
