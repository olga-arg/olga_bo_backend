package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-delete-team/internal/processor"
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
	fmt.Println("Extracting team ID from URL path parameter")
	teamID, ok := request.PathParameters["team_id"]
	if !ok {
		err := errors.New("missing team ID in request")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, err
	}

	// Update team in storage
	fmt.Println("Updating team in storage")
	err := h.processor.DeleteTeam(context.Background(), teamID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	// Convert team to DTO and write response
	fmt.Println("Converting team to DTO and writing response")
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Team deleted successfully",
	}, nil
}
