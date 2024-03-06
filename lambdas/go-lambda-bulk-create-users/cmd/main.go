package main

import (
	"commons/utils/db"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-bulk-create-users/internal/application"
	"go-lambda-bulk-create-users/internal/processor"
	"go-lambda-bulk-create-users/pkg/handler"
)

func main() {
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
