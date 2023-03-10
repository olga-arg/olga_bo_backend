package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"go-lambda-get-all-users/internal/processor"
)

type GetAllUsersHandler struct {
	processor processor.Processor
}

func NewGetAllUsersHandler(p processor.Processor) *GetAllUsersHandler {
	return &GetAllUsersHandler{processor: p}
}

func (h *GetAllUsersHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	users, err := h.processor.GetAllUsers(context.Background())
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	responseBody := map[string]interface{}{
		"users": users,
	}

	body, err := json.Marshal(responseBody)
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
