package handler

import (
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
	var input dto.CreateCompanyInput

	email, companyId, err := utils.ExtractEmailAndCompanyIdFromToken(request)
	if err != nil {
		fmt.Println("Error extracting email and company id from token: ", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       err.Error(),
		}, nil
	}

	if companyId == "" || email == "" {
		println("companyId or email is empty", companyId, email)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unauthorized",
		}, nil
	}

	if "vilavalentin@gmail.com" != email && "iramosibx@gmail.com" != email {
		println("email is not admin", email)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unauthorized",
		}, nil
	}

	// Validate input
	err = h.processor.ValidateCompanyInput(context.Background(), &input, request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	err = h.processor.CreateCompany(context.Background(), &input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	err = services.NewDefaultEmailService().SendEmail(input.UserEmail, services.Welcome, []string{input.UserName}, nil)
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
		Body: "Company created successfully",
	}, nil
}
