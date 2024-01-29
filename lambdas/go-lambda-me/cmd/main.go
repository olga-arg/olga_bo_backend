package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-me/internal/application"
	"go-lambda-me/internal/processor"
	"go-lambda-me/internal/storage"
	"go-lambda-me/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	db, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	paymentRepo := storage.NewPaymentRepository(db)
	paymentProcessor := processor.NewProcessor(paymentRepo)
	getAllPaymentsHandler := handler.NewMeHandler(paymentProcessor)
	lambda.Start(getAllPaymentsHandler.Handle)
}
