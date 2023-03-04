package handler

import (
	"context"
	"encoding/json"
	"go-lambda-create-user/internal/application"
	"go-lambda-create-user/pkg/dto"
	"net/http"
)

func CreateUser(app *application.App) func(ctx context.Context, request *http.Request) (response interface{}, err error) {
	return func(ctx context.Context, request *http.Request) (response interface{}, err error) {
		var input dto.UserInput
		if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
			return nil, err
		}
		userOutput, err := app.Processor.CreateUser(ctx, &input)
		if err != nil {
			return nil, err
		}
		response, err = json.Marshal(userOutput)
		if err != nil {
			return nil, err
		}
		return response, nil
	}
}
