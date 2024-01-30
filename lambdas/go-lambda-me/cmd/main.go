package main

import (
	"commons/utils/db"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-me/internal/application"
	"go-lambda-me/internal/processor"
	"go-lambda-me/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	gormDb, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	paymentRepo := db.NewPaymentRepository(gormDb)
	userRepo := db.NewUserRepository(gormDb)
	paymentProcessor := processor.NewProcessor(*paymentRepo, *userRepo)
	getAllPaymentsHandler := handler.NewMeHandler(paymentProcessor)
	lambda.Start(getAllPaymentsHandler.Handle)
}
