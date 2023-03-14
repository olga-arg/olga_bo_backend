package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-update-team-annual-budget/internal/processor"
	"net/http"
	"strconv"
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
	// Get the annual budget from the request body
	fmt.Println("Extracting annual budget from request body")
	annualBudget, ok := request.QueryStringParameters["annual_budget"]
	annualBudgetInt, err := strconv.Atoi(annualBudget)

	// Update team in storage
	fmt.Println("Updating team in storage")
	err = h.processor.UpdateTeamBudget(context.Background(), teamID, annualBudgetInt)
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
		Body:       "Team Updated successfully",
	}, nil
}
