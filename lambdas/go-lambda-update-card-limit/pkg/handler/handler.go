package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-update-card-limit/internal/processor"
	"go-lambda-update-card-limit/pkg/domain"
	"go-lambda-update-card-limit/pkg/dto"
	"net/http"
	"time"
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

	var input struct {
		PurchaseLimit int       `json:"purchase_limit"`
		MonthlyLimit  int       `json:"monthly_limit"`
		ResetDate     time.Time `json:"reset_date"`
	}
	if err := json.Unmarshal([]byte(request.Body), &input); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "invalid request body",
		}, nil
	}

	// Validate input
	if input.PurchaseLimit < 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "invalid purchase limit",
		}, nil
	}
	if input.MonthlyLimit < 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "invalid monthly limit",
		}, nil
	}
	if input.ResetDate.IsZero() {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "invalid reset date",
		}, nil
	}

	// Update user in storage
	err := h.processor.UpdateUserCardLimits(context.Background(), userID, input.PurchaseLimit, input.MonthlyLimit)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "failed to update user in storage",
		}, nil
	}

	// Update reset date in storage
	err = h.processor.UpdateUserResetDate(context.Background(), userID, input.ResetDate)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "failed to update reset date in storage",
		}, nil
	}

	// Get updated user from storage
	updatedUser, err := h.processor.GetUser(context.Background(), userID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "failed to get updated user from storage",
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
