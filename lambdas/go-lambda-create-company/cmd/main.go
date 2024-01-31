package main

import (
	"commons/utils/db"
	"github.com/aws/aws-lambda-go/lambda"
	"go-lambda-create-user/internal/application"
	"go-lambda-create-user/internal/processor"
	"go-lambda-create-user/pkg/handler"
)

func main() {
	application.SetupEmailService()
	pgConnector := application.PostgresConnector{}
	gormDb, err := pgConnector.GetConnection()
	if err != nil {
		panic(err)
	}
	companyRepo := db.NewCompanyRepository(gormDb)
	userRepo := db.NewUserRepository(gormDb)
	categoryRepo := db.NewCategoryRepository(gormDb)
	userProcessor := processor.New(*companyRepo, *userRepo, *categoryRepo)
	createUserHandler := handler.NewCreateUserHandler(userProcessor)
	lambda.Start(createUserHandler.Handle)
}
