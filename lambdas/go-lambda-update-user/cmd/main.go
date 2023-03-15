package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-update-user/internal/application"
	"go-lambda-update-user/internal/processor"
	"go-lambda-update-user/internal/storage"
	"go-lambda-update-user/pkg/handler"
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
