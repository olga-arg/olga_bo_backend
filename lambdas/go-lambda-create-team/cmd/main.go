package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-create-team/internal/application"
	"go-lambda-create-team/internal/processor"
	"go-lambda-create-team/internal/storage"
	"go-lambda-create-team/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	db, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	teamRepo := storage.NewUserRepository(db)
	teamProcessor := processor.New(*teamRepo)
	createTeamHandler := handler.NewCreateTeamHandler(teamProcessor)
	lambda.Start(createTeamHandler.Handle)
}
