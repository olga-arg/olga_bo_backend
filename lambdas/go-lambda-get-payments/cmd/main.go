package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-get-payments/internal/application"
	"go-lambda-get-payments/internal/processor"
	"go-lambda-get-payments/internal/storage"
	"go-lambda-get-payments/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	db, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	paymentRepo := storage.NewPaymentRepository(db)
	paymentProcessor := processor.NewProcessor(paymentRepo)
	getAllPaymentsHandler := handler.NewGetAllPaymentsHandler(paymentProcessor)
	lambda.Start(getAllPaymentsHandler.Handle)
}
