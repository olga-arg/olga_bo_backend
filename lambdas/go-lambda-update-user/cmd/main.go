package main

import (
	"commons/utils/db"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-update-user/internal/application"
	"go-lambda-update-user/internal/processor"
	"go-lambda-update-user/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	gormDb, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	userRepo := db.NewUserRepository(gormDb)
	userProcessor := processor.NewProcessor(userRepo)
	getAllUsersHandler := handler.NewUserCardLimitHandler(userProcessor)
	lambda.Start(getAllUsersHandler.Handle)
}
