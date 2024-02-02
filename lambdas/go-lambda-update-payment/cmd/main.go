package main

import (
	"commons/utils/db"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-update-payment/internal/application"
	"go-lambda-update-payment/internal/processor"
	"go-lambda-update-payment/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	gormDb, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	userRepo := db.NewUserRepository(gormDb)
	teamRepo := db.NewTeamRepository(gormDb)
	paymentRepo := db.NewPaymentRepository(gormDb)
	paymentProcessor := processor.NewProcessor(*paymentRepo, *teamRepo, *userRepo)
	getAllPaymentsHandler := handler.NewUpdatePaymentHandler(paymentProcessor)
	lambda.Start(getAllPaymentsHandler.Handle)
}
