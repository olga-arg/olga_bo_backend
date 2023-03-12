package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-get-all-users-rds/internal/application"
	"go-lambda-get-all-users-rds/internal/processor"
	"go-lambda-get-all-users-rds/internal/storage"
	"go-lambda-get-all-users-rds/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	db, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	userRepo := storage.NewUserRepository(db)
	userProcessor := processor.NewProcessor(userRepo)
	getAllUsersHandler := handler.NewGetAllUsersHandler(userProcessor)
	lambda.Start(getAllUsersHandler.Handle)
}
