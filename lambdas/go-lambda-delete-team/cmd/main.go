package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-delete-team/internal/application"
	"go-lambda-delete-team/internal/processor"
	"go-lambda-delete-team/internal/storage"
	"go-lambda-delete-team/pkg/handler"
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
