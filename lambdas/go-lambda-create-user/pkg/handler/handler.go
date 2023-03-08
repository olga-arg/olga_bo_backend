package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-create-user/internal/processor"
	"go-lambda-create-user/pkg/dto"
	"go-lambda-create-user/internal/services"
	"net/http"
	"github.com/go-playground/validator/v10"
)

type CreateUserHandler struct {
	processor processor.Processor
}

func NewCreateUserHandler(p processor.Processor) *CreateUserHandler {
	return &CreateUserHandler{processor: p}
}

var validate = validator.New()

func (h *CreateUserHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.Body == "" || len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing request body",
		}, nil
	}

	// Creates a CreateUserInput struct from the request body
	var input dto.CreateUserInput
	// Unmarshal the request body into the CreateUserInput struct
	err := json.Unmarshal([]byte(request.Body), &input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, nil
	}

	// Validates that the JSON request body has the correct fields and that they are the correct type
	if err := validate.Struct(input); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing or invalid fields in request body",
		}, nil
	}

	output, err := h.processor.CreateUser(context.Background(), &input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	// Send email to user
	if err := services.SendEmail(
		"Welcome to the team!",
		"You have been added to the team. Please log in to your account to view your teams.",
		[]string{input.Email},
		nil, nil, nil,
		)

	responseBody, err := json.Marshal(output)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "User created successfully",
	}, nil
}

