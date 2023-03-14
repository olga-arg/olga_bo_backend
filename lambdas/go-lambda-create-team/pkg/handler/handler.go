package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-create-team/internal/processor"
	"go-lambda-create-team/pkg/dto"
	"net/http"
)

type CreateTeamHandler struct {
	processor processor.Processor
}

func NewCreateTeamHandler(p processor.Processor) *CreateTeamHandler {
	return &CreateTeamHandler{processor: p}
}

func (h *CreateTeamHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.Body == "" || len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing request body",
		}, nil
	}

	var input dto.CreateTeamInput
	err := json.Unmarshal([]byte(request.Body), &input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
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
