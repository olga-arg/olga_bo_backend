package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-create-user/internal/application"
	"go-lambda-create-user/internal/processor"
	"go-lambda-create-user/internal/storage"
	"go-lambda-create-user/pkg/handler"
)

func main() {
	application.SetupEmailService()
	pgConnector := application.PostgresConnector{}
	db, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	userRepo := storage.NewUserRepository(db)
	userProcessor := processor.New(*userRepo)
	createUserHandler := handler.NewCreateUserHandler(userProcessor)
	lambda.Start(createUserHandler.Handle)
}
