package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-update-payment/internal/application"
	"go-lambda-update-payment/internal/processor"
	"go-lambda-update-payment/internal/storage"
	"go-lambda-update-payment/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	db, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	userRepo := storage.NewPaymentRepository(db)
	paymentProcessor := processor.NewProcessor(userRepo)
	getAllPaymentsHandler := handler.NewUpdatePaymentHandler(paymentProcessor)
	lambda.Start(getAllPaymentsHandler.Handle)
}
