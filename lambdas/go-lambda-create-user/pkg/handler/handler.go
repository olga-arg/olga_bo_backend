package handler

import (
	"context"
	"encoding/json"
	"go-lambda-create-user/internal/application"
	"go-lambda-create-user/pkg/dto"
	"log"
	"net/http"
)

func CreateUser(app *application.App) func(ctx context.Context, request *http.Request) (response interface{}, err error) {
	return func(ctx context.Context, request *http.Request) (response interface{}, err error) {
		var input dto.UserInput
		log.Println("input: ", input)
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			log.Println("error 1: ")
			return nil, err
		}
		userOutput, err := app.Processor.CreateUser(ctx, &input)
		log.Println("userOutput: ", userOutput)
		if err != nil {
			log.Println("error 2: ")
			return nil, err
		}
		log.Println("going to be marshalled: ")
		response, err = json.Marshal(userOutput)
		log.Println("response: ", response)
		if err != nil {
			log.Println("error 3: ")
			return nil, err
		}
		return response, nil
	}
}
