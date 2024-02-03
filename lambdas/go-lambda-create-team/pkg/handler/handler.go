package handler

import (
	"commons/domain"
	"commons/utils"
	"context"
	"fmt"
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
	var input dto.CreateTeamInput

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

	allowedRoles := []domain.UserRoles{domain.Admin}
	isAuthorized, err := h.processor.ValidateUser(context.Background(), email, companyId, allowedRoles)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       err.Error(),
		}, nil
	}
	if !isAuthorized {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unauthorized",
		}, nil
	}

	// Validate input
	fmt.Println("Validating input")
	err = h.processor.ValidateTeamInput(context.Background(), &input, request, companyId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	err = h.processor.CreateTeam(context.Background(), &input, companyId)
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
