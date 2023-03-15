package handler

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-create-payment/internal/processor"
	"net/http"
)

type CreatePaymentHandler struct {
	processor processor.Processor
}

func NewCreatePaymentHandler(p processor.Processor) *CreatePaymentHandler {
	return &CreatePaymentHandler{processor: p}
}

func (h *CreatePaymentHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	err := h.processor.CreatePayment(context.Background())
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
