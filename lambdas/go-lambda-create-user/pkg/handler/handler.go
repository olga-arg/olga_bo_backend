package handler

import (
	"commons/domain"
	"commons/utils"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-create-user/internal/processor"
	"go-lambda-create-user/internal/services"
	"go-lambda-create-user/pkg/dto"
	"net/http"
)

type CreateUserHandler struct {
	processor processor.Processor
}

func NewCreateUserHandler(p processor.Processor) *CreateUserHandler {
	return &CreateUserHandler{processor: p}
}

func (h *CreateUserHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input dto.CreateUserInput

	email, companyId, err := utils.ExtractEmailAndCompanyIdFromToken(request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       err.Error(),
		}, nil
	}

	if companyId == "" || email == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unauthorized",
		}, nil
	}

	// Validate user
	allowedRoles := []domain.UserRoles{domain.Admin}
	isAuthorized, err := h.processor.ValidateUser(context.Background(), email, companyId, allowedRoles)
	if err != nil {
		fmt.Println("Error validating user: ", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       err.Error(),
		}, nil
	}
	if !isAuthorized {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unauthorized",
		}, nil
	}

	// Validate input
	err = h.processor.ValidateUserInput(context.Background(), &input, request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	err = h.processor.CreateUser(context.Background(), &input, companyId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	err = services.NewDefaultEmailService().SendEmail(input.Email, services.Welcome, []string{input.Name}, nil)
	if err != nil {
		fmt.Println("Error sending email:", err)
	} else {
		fmt.Println("Email sent successfully.")
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
		},
		Body: "User created successfully",
	}, nil
}
