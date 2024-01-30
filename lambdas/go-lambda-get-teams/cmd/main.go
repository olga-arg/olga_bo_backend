package main

import (
	"commons/utils/db"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-get-teams/internal/application"
	"go-lambda-get-teams/internal/processor"
	"go-lambda-get-teams/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	gormDb, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	teamRepo := db.NewTeamRepository(gormDb)
	teamProcessor := processor.NewProcessor(teamRepo)
	getAllTeamsHandler := handler.NewGetAllTeamsHandler(teamProcessor)
	lambda.Start(getAllTeamsHandler.Handle)
}
