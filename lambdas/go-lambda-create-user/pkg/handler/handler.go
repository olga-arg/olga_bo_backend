package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"go-lambda-create-user/internal/processor"
	"go-lambda-create-user/internal/services"
	"go-lambda-create-user/pkg/dto"
	"net/http"
	"os"
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

	_, err = h.processor.CreateUser(context.Background(), &input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	// Send email to user
	err = godotenv.Load("../../../../.env")
	fromEmailAddress := os.Getenv("EMAIL_SENDER_ADDRESS")
	fromEmailPassword := os.Getenv("EMAIL_SENDER_PASSWORD")
	sender := services.NewEmailSender(fromEmailAddress, fromEmailPassword)
	subject := "Test email"
	body := "This is a test email"
	to := []string{"ir.basura@gmail.com"}

	err = sender.SendEmail(subject, body, to, nil, nil, nil)

	// responseBody, _ := json.Marshal(output)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "User created successfully",
	}, nil
}
