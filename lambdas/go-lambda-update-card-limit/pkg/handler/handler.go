package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-update-card-limit/internal/processor"
	"go-lambda-update-card-limit/pkg/dto"

	"net/http"
)

type UserCardLimitHandler struct {
	processor processor.Processor
}

func NewUserCardLimitHandler(processor processor.Processor) *UserCardLimitHandler {
	return &UserCardLimitHandler{
		processor: processor,
	}
}

func (h *UserCardLimitHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input dto.UpdateLimitInput

	// Validate input
	fmt.Println("Validating input")
	newUser, err := h.processor.ValidateUserInput(context.Background(), &input, request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}
	// Update user in storage
	fmt.Println("Updating user in storage")
	err = h.processor.UpdateUserCardLimits(context.Background(), newUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "failed to update user in storage",
		}, nil
	}

	// Convert user to DTO and write response
	fmt.Println("Converting user to DTO and writing response")
	output := dto.NewOutput(newUser)
	fmt.Println("output:", output)
	responseBody, err := json.Marshal(output)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "failed to encode response",
		}, nil
	}
	fmt.Println("response:", string(responseBody))
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
	}, nil
}
