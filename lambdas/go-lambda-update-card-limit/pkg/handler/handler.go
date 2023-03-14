package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-update-card-limit/internal/processor"
	"go-lambda-update-card-limit/pkg/domain"
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
	userID := request.PathParameters["user_id"]

	var input dto.UpdateLimitInput

	// Validate input
	if err := h.processor.ValidateUserInput(context.Background(), &input, request); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	// Update user in storage
	updatedUser, err := h.processor.UpdateUserCardLimits(context.Background(), userID, input.PurchaseLimit, input.MonthlyLimit)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "failed to update user in storage",
		}, nil
	}

	// Convert user to DTO and write response
	output := dto.NewOutput([]domain.User{*updatedUser})
	responseBody, err := json.Marshal(output)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "failed to encode response",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBody),
	}, nil
}
