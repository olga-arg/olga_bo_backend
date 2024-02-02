package main

import (
	"commons/utils/db"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-create-payment/internal/application"
	"go-lambda-create-payment/internal/processor"
	"go-lambda-create-payment/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	gormDb, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	paymentRepo := db.NewPaymentRepository(gormDb)
	userRepo := db.NewUserRepository(gormDb)
	teamRepo := db.NewTeamRepository(gormDb)
	paymentProcessor := processor.New(*paymentRepo, *userRepo, *teamRepo)
	createPaymentHandler := handler.NewCreatePaymentHandler(paymentProcessor)
	lambda.Start(createPaymentHandler.Handle)
}
