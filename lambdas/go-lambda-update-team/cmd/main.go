package main

import (
	"commons/utils/db"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-update-team/internal/application"
	"go-lambda-update-team/internal/processor"
	"go-lambda-update-team/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	gormDb, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	teamRepo := db.NewTeamRepository(gormDb)
	teamProcessor := processor.NewProcessor(teamRepo)
	deleteTeamHandler := handler.NewTeamHandler(teamProcessor)
	lambda.Start(deleteTeamHandler.Handle)
}
