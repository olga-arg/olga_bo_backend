package handler

import (
	"commons/domain"
	"commons/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-update-team/internal/processor"
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
	email, companyId, err := utils.ExtractEmailAndCompanyIdFromToken(request)
	if err != nil {
		fmt.Println("Error extracting email and company id from token: ", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       err.Error(),
		}, nil
	}

	if companyId == "" || email == "" {
		println("companyId or email is empty", companyId, email)
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

	// Parse the request body into a struct
	fmt.Println("Parsing request body")
	var updateRequest *domain.UpdateTeamRequest
	err = json.Unmarshal([]byte(request.Body), &updateRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, err
	}

	// Update team in storage
	fmt.Println("Updating team in storage")
	err = h.processor.UpdateTeam(context.Background(), teamID, updateRequest, companyId)
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
