package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-get-all-teams/internal/application"
	"go-lambda-get-all-teams/internal/processor"
	"go-lambda-get-all-teams/internal/storage"
	"go-lambda-get-all-teams/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	db, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	teamRepo := storage.NewTeamRepository(db)
	teamProcessor := processor.NewProcessor(teamRepo)
	getAllTeamsHandler := handler.NewGetAllTeamsHandler(teamProcessor)
	lambda.Start(getAllTeamsHandler.Handle)
}
