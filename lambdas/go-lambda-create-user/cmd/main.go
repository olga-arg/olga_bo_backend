package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-create-user/internal/application"
	"go-lambda-create-user/pkg/handler"
)

func main() {
	ctx := context.Background()
	app, err := application.New(ctx)
	if err != nil {
		panic(err)
	}
	lambda.Start(handler.CreateUser(app))
}
