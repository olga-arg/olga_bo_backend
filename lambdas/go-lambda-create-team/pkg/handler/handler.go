package handler

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-create-team/internal/processor"
	"go-lambda-create-team/pkg/dto"
	"log"
	"net/http"
)

type CreateTeamHandler struct {
	processor processor.Processor
}

func NewCreateTeamHandler(p processor.Processor) *CreateTeamHandler {
	return &CreateTeamHandler{processor: p}
}

func (h *CreateTeamHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input dto.CreateTeamInput

	// Validate input
	log.Println("Validating input")
	err := h.processor.ValidateTeamInput(context.Background(), &input, request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	err = h.processor.CreateTeam(context.Background(), &input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "Team created successfully",
	}, nil
}
