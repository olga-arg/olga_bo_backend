package handler

import (
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
	var input dto.UpdatePaymentInput

	// Validate input
	fmt.Println("Validating input")
	newPayment, err := h.processor.ValidatePaymentInput(context.Background(), &input, request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}
	// Update user in storage
	fmt.Println("Updating payment in storage")
	err = h.processor.UpdatePayment(context.Background(), newPayment)
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
