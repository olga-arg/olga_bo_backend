package main

import (
	"commons/utils/db"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-post-confirmation/internal/application"
	"go-lambda-post-confirmation/internal/processor"
	"go-lambda-post-confirmation/pkg/handler"
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
