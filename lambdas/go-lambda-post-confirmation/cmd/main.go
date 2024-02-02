package main

import (
	"commons/utils/db"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-get-all-users/internal/application"
	"go-lambda-get-all-users/internal/processor"
	"go-lambda-get-all-users/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	gormDb, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	userRepo := db.NewUserRepository(gormDb)
	userProcessor := processor.NewProcessor(userRepo)
	postConfirmationHandler := handler.NewPostConfirmationHandler(userProcessor)
	lambda.Start(postConfirmationHandler.Handle)
}
