package handler

import (
	"commons/utils"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-get-teams/internal/processor"
	"net/http"
)

type GetAllTeamsHandler struct {
	processor processor.Processor
}

func NewGetAllTeamsHandler(p processor.Processor) *GetAllTeamsHandler {
	return &GetAllTeamsHandler{processor: p}
}

func (h *GetAllTeamsHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	filters := request.QueryStringParameters

	teams, err := h.processor.GetAllTeams(context.Background(), filters, companyId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}
	body, err := json.Marshal(teams)
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
