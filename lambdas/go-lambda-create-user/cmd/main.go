package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-create-user/pkg/handler"
)

func main() {
	lambda.Start(handler.CreateUser)
}
