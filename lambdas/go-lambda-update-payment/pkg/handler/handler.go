package handler

import (
	"commons/domain"
	"commons/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-update-payment/internal/processor"
	"go-lambda-update-payment/pkg/dto"

	"net/http"
)

type UpdatePaymentHandler struct {
	processor processor.Processor
}

func NewUpdatePaymentHandler(processor processor.Processor) *UpdatePaymentHandler {
	return &UpdatePaymentHandler{
		processor: processor,
	}
}

func (h *UpdatePaymentHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	allowedRoles := []domain.UserRoles{domain.Reviewer, domain.Admin}
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

	var input dto.UpdatePaymentInput

	// Validate input
	fmt.Println("Validating input")
	newPayment, err := h.processor.ValidatePaymentInput(context.Background(), &input, request, companyId, email)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}
	// Update user in storage
	fmt.Println("Updating payment in storage")
	err = h.processor.UpdatePayment(context.Background(), newPayment, companyId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "failed to update payment in storage",
		}, nil
	}

	// Convert user to DTO and write response
	fmt.Println("Converting payment to DTO and writing response")
	output := dto.NewOutput(newPayment)
	fmt.Println("output:", output)
	responseBody, err := json.Marshal(output)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "failed to encode response",
		}, nil
	}
	fmt.Println("response:", string(responseBody))
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
	}, nil
}
