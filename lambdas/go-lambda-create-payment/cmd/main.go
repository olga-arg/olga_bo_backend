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
	paymentProcessor := processor.New(*paymentRepo)
	createPaymentHandler := handler.NewCreatePaymentHandler(paymentProcessor)
	lambda.Start(createPaymentHandler.Handle)
}
