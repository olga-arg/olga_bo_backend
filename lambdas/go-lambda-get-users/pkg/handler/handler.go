package handler

import (
	"commons/domain"
	"commons/utils"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-get-all-users/internal/processor"
	"net/http"
)

type GetAllUsersHandler struct {
	processor processor.Processor
}

func NewGetAllUsersHandler(p processor.Processor) *GetAllUsersHandler {
	return &GetAllUsersHandler{processor: p}
}

func (h *GetAllUsersHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	allowedRoles := []domain.UserRoles{domain.Reviewer, domain.Admin, domain.Accountant}
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

	filters := request.QueryStringParameters

	users, err := h.processor.GetAllUsers(context.Background(), filters, companyId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}
	body, err := json.Marshal(users) // TODO: Pasar users a un Response DTO, por si quiero devolver menos cosas.
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(body),
	}, nil
}
