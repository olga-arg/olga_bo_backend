package handler

import (
	"context"
	"fmt"
	"go-lambda-create-user/internal/application"
	"net/http"
)

func CreateUser(app *application.App) func(ctx context.Context, request *http.Request) (response interface{}, err error) {
	fmt.Println("TESTING!!!")
	return func(ctx context.Context, request *http.Request) (response interface{}, err error) {
		fmt.Println("TESTING2!!!")
		return response, err
	}
}
