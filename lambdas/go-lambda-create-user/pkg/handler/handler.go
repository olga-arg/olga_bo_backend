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
	return func(ctx context.Context, request *http.Request) (response interface{}, err error) {
		var input dto.UserInput
		fmt.Println("input: ", input)
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			fmt.Println("error 1: ")
			return nil, err
		}
		userOutput, err := app.Processor.CreateUser(ctx, &input)
		fmt.Println("userOutput: ", userOutput)
		if err != nil {
			fmt.Println("error 2: ")
			return nil, err
		}
		fmt.Println("going to be marshalled: ")
		response, err = json.Marshal(userOutput)
		fmt.Println("response: ", response)
		if err != nil {
			fmt.Println("error 3: ")
			return nil, err
		}
		return response, nil
	}
}
