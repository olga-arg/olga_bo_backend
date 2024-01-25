package main

import (
	"commons/utils/db"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-create-team/internal/application"
	"go-lambda-create-team/internal/processor"
	"go-lambda-create-team/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	gormDb, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	teamRepo := db.NewTeamRepository(gormDb)
	teamProcessor := processor.New(*teamRepo)
	createTeamHandler := handler.NewCreateTeamHandler(teamProcessor)
	lambda.Start(createTeamHandler.Handle)
}
