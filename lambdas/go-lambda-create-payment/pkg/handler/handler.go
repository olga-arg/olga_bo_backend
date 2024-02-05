package handler

import (
	"commons/domain"
	"commons/utils"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-create-payment/internal/processor"
	"go-lambda-create-payment/pkg/dto"
	"net/http"
)

type CreatePaymentHandler struct {
	processor processor.Processor
}

func NewCreatePaymentHandler(p processor.Processor) *CreatePaymentHandler {
	return &CreatePaymentHandler{processor: p}
}

func (h *CreatePaymentHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input dto.CreatePaymentInput
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
	allowedRoles := []domain.UserRoles{domain.Employee, domain.Reviewer, domain.Admin}
	isAuthorized, err := h.processor.ValidateUser(context.Background(), email, companyId, allowedRoles)
	if err != nil {
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

	// Unmarshal the request body into the input struct
	err = json.Unmarshal([]byte(request.Body), &input)
	if err != nil {
		print(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	err = h.processor.CreatePayment(context.Background(), &input, email, companyId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "Payment created successfully",
	}, nil
}
