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
	userRepo := db.NewUserRepository(gormDb)
	teamProcessor := processor.NewProcessor(*teamRepo, *userRepo)
	deleteTeamHandler := handler.NewTeamHandler(teamProcessor)
	lambda.Start(deleteTeamHandler.Handle)
}
