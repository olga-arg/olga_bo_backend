package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-delete-team/internal/processor"
	"go-lambda-delete-team/pkg/dto"
	"log"
	"net/http"
)

type TeamHandler struct {
	processor processor.Processor
}

func NewTeamHandler(processor processor.Processor) *TeamHandler {
	return &TeamHandler{
		processor: processor,
	}
}

func (h *TeamHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Extract team ID from URL path parameter
	teamID, ok := request.PathParameters["team_id"]
	if !ok {
		err := errors.New("missing team ID in request")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, err
	}

	team, err := h.processor.GetTeam(context.Background(), teamID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, err
	}

	if team == nil {
		err := errors.New("team not found")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       err.Error(),
		}, err
	}

	// Update team in storage
	log.Println("Updating team in storage")
	err = h.processor.DeleteTeam(context.Background(), team)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "failed to update team in storage",
		}, nil
	}

	// Convert team to DTO and write response
	log.Println("Converting team to DTO and writing response")
	output := dto.NewOutput(team)
	log.Println("output:", output)
	responseBody, err := json.Marshal(output)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "failed to encode response",
		}, nil
	}
	log.Println("response:", string(responseBody))
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
	}, nil
}
