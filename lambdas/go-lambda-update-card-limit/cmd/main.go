package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-update-card-limit/internal/application"
	"go-lambda-update-card-limit/internal/processor"
	"go-lambda-update-card-limit/internal/storage"
	"go-lambda-update-card-limit/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	db, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	userRepo := storage.NewUserRepository(db)
	userProcessor := processor.NewProcessor(userRepo)
	getAllUsersHandler := handler.NewUserCardLimitHandler(userProcessor)
	lambda.Start(getAllUsersHandler.Handle)
}
