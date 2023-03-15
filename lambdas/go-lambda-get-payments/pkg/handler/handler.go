package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-get-payments/internal/processor"
	"net/http"
)

type GetAllPaymentsHandler struct {
	processor processor.Processor
}

func NewGetAllPaymentsHandler(p processor.Processor) *GetAllPaymentsHandler {
	return &GetAllPaymentsHandler{processor: p}
}

func (h *GetAllPaymentsHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	filters := request.QueryStringParameters

	users, err := h.processor.GetAllPayments(context.Background(), filters)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}
	body, err := json.Marshal(users)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}
