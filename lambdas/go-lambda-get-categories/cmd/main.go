package main

import (
	"commons/utils/db"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-get-categories/internal/application"
	"go-lambda-get-categories/internal/processor"
	"go-lambda-get-categories/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	gormDb, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	categoryRepo := db.NewCategoryRepository(gormDb)
	userRepo := db.NewUserRepository(gormDb)
	categoryProcessor := processor.NewProcessor(*categoryRepo, *userRepo)
	getCategoriesHandler := handler.NewGetCategoriesHandler(categoryProcessor)
	lambda.Start(getCategoriesHandler.Handle)
}
