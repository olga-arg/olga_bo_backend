package main

import (
	"commons/utils/db"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-export-payments/internal/application"
	"go-lambda-export-payments/internal/processor"
	"go-lambda-export-payments/pkg/handler"
)

func main() {
	pgConnector := application.PostgresConnector{}
	gormDb, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	paymentRepo := db.NewPaymentRepository(gormDb)
	paymentProcessor := processor.New(*paymentRepo)
	createUserHandler := handler.NewExportPaymentHandler(paymentProcessor)
	lambda.Start(createUserHandler.Handle)
}
