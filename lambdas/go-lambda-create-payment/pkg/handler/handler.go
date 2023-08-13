package handler

import (
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

	// Unmarshal the request body into the input struct
	err := json.Unmarshal([]byte(request.Body), &input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	err = h.processor.CreatePayment(context.Background(), &input)
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
