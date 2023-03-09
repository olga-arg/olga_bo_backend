package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-get-all-users/internal/application"
	"go-lambda-get-all-users/internal/processor"
	"go-lambda-get-all-users/internal/storage"
	"go-lambda-get-all-users/pkg/handler"
)

func main() {
	db := application.NewDynamoDBClient()
	userRepo := storage.NewUserRepository(db)
	userProcessor := processor.NewProcessor(userRepo)
	getAllUsersHandler := handler.NewGetAllUsersHandler(userProcessor)
	lambda.Start(getAllUsersHandler.Handle)
}
