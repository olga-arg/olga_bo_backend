package main

import (
	"commons/utils/db"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-create-user/internal/application"
	"go-lambda-create-user/internal/processor"
	"go-lambda-create-user/pkg/handler"
)

func main() {
	application.SetupEmailService()
	pgConnector := application.PostgresConnector{}
	gormDb, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	userRepo := db.NewUserRepository(gormDb)
	userProcessor := processor.New(*userRepo)
	createUserHandler := handler.NewCreateUserHandler(userProcessor)
	lambda.Start(createUserHandler.Handle)
}
