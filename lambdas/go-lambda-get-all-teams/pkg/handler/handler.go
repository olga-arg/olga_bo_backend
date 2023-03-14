package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-get-all-teams/internal/processor"
	"net/http"
)

type GetAllTeamsHandler struct {
	processor processor.Processor
}

func NewGetAllTeamsHandler(p processor.Processor) *GetAllTeamsHandler {
	return &GetAllTeamsHandler{processor: p}
}

func (h *GetAllTeamsHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	filters := request.QueryStringParameters

	users, err := h.processor.GetAllTeams(context.Background(), filters)
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
