package main

import (
	"commons/utils/db"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-create-category/internal/application"
	"go-lambda-create-category/internal/processor"
	"go-lambda-create-category/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	gormDb, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	categoryRepo := db.NewCategoryRepository(gormDb)
	userRepo := db.NewUserRepository(gormDb)
	categoryProcessor := processor.New(*categoryRepo, *userRepo)
	createCategoryHandler := handler.NewCreateCategoryHandler(categoryProcessor)
	lambda.Start(createCategoryHandler.Handle)
}
