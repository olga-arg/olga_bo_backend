package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-create-payment/internal/application"
	"go-lambda-create-payment/internal/processor"
	"go-lambda-create-payment/internal/storage"
	"go-lambda-create-payment/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	db, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	paymentRepo := storage.NewPaymentRepository(db)
	paymentProcessor := processor.New(*paymentRepo)
	createPaymentHandler := handler.NewCreatePaymentHandler(paymentProcessor)
	lambda.Start(createPaymentHandler.Handle)
}
