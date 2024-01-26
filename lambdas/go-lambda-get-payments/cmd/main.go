package main

import (
	"commons/utils/db"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-get-payments/internal/application"
	"go-lambda-get-payments/internal/processor"
	"go-lambda-get-payments/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	gormDb, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	paymentRepo := db.NewPaymentRepository(gormDb)
	paymentProcessor := processor.NewProcessor(paymentRepo)
	getAllPaymentsHandler := handler.NewGetAllPaymentsHandler(paymentProcessor)
	lambda.Start(getAllPaymentsHandler.Handle)
}
