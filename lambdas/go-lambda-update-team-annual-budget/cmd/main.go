package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-update-team-annual-budget/internal/application"
	"go-lambda-update-team-annual-budget/internal/processor"
	"go-lambda-update-team-annual-budget/internal/storage"
	"go-lambda-update-team-annual-budget/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	db, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	teamRepo := storage.NewTeamRepository(db)
	teamProcessor := processor.NewProcessor(teamRepo)
	deleteTeamHandler := handler.NewTeamHandler(teamProcessor)
	lambda.Start(deleteTeamHandler.Handle)
}
