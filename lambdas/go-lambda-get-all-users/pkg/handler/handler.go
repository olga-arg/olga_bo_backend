package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-get-all-users/internal/processor"
	"net/http"
)

type GetAllUsersHandler struct {
	processor processor.Processor
}

func NewGetAllUsersHandler(p processor.Processor) *GetAllUsersHandler {
	return &GetAllUsersHandler{processor: p}
}

func (h *GetAllUsersHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	filter := request.QueryStringParameters["filter"]
	users, err := h.processor.GetAllUsers(context.Background(), filter)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	response := map[string]interface{}{
		"users":  users.Users,
		"filter": filter,
	}

	body, err := json.Marshal(response)
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
