package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/go-playground/validator/v10"
	"go-lambda-create-user/internal/processor"
	"go-lambda-create-user/internal/services"
	"go-lambda-create-user/pkg/dto"
	"net/http"
)

type CreateUserHandler struct {
	processor processor.Processor
	ssmClient ssmiface.SSMAPI
}

func NewCreateUserHandler(p processor.Processor) *CreateUserHandler {
	// Creates a new session using the default AWS configuration
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Creates a new SSM parameter store interface using the session
	ssmClient := ssm.New(sess)

	return &CreateUserHandler{
		processor: p,
		ssmClient: ssmClient,
	}
}

var validate = validator.New()

func (h *CreateUserHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fromEmailAddressParam, err := h.ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String("/olga-backend/EMAIL_SENDER_ADDRESS"),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error retrieving email sender address parameter from SSM",
		}, nil
	}
	fromEmailPasswordParam, err := h.ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String("/olga-backend/EMAIL_SENDER_PASSWORD"),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error retrieving email sender password parameter from SSM",
		}, nil
	}

	// Get the value of the decrypted email sender address and password from the SSM parameter store
	fromEmailAddress := aws.StringValue(fromEmailAddressParam.Parameter.Value)
	fromEmailPassword := aws.StringValue(fromEmailPasswordParam.Parameter.Value)

	if request.Body == "" || len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing request body",
		}, nil
	}

	// Creates a CreateUserInput struct from the request body
	var input dto.CreateUserInput
	// Unmarshal the request body into the CreateUserInput struct
	err = json.Unmarshal([]byte(request.Body), &input)
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

	// Send email to user
	sender := services.NewEmailSender(fromEmailAddress, fromEmailPassword)
	subject := "Test email"
	body := "This is a test email"
	to := []string{input.Email}

	err = sender.SendEmail(subject, body, to, nil, nil, nil)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Error sending email from %s to %s", fromEmailAddress, input.Email),
		}, nil
	}

	_, err = h.processor.CreateUser(context.Background(), &input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	// responseBody, _ := json.Marshal(output)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "User created successfully",
	}, nil
}
