package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"go-lambda-create-user/internal/application"
	"go-lambda-create-user/pkg/dto"
	"net/http"
)

func CreateUser(app *application.App) func(ctx context.Context, request *http.Request) (response interface{}, err error) {
	fmt.Println("TESTING2")
	return func(ctx context.Context, request *http.Request) (response interface{}, err error) {
		var input dto.UserInput
		fmt.Println("input: ", input)
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			fmt.Println("error 1: ")
			return nil, err
		}
		fmt.Println("BUENAS")
		return processCreateUserRequest(ctx, app, &input)
	}
}

func processCreateUserRequest(ctx context.Context, app *application.App, input *dto.UserInput) (response interface{}, err error) {
	fmt.Println("BUENAS2")
	userOutput, err := app.Processor.CreateUser(ctx, input)
	if err != nil {
		return nil, err
	}

	response, err = json.Marshal(userOutput)
	if err != nil {
		return nil, err
	}

	return response, nil
}
