package handler

import (
	"commons/utils"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-get-categories/internal/processor"
	"net/http"
)

type GetCategoriesHandler struct {
	processor processor.Processor
}

func NewGetCategoriesHandler(p processor.Processor) *GetCategoriesHandler {
	return &GetCategoriesHandler{processor: p}
}

func (h *GetCategoriesHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	categories, err := h.processor.GetCategories(context.Background(), companyId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}
	body, err := json.Marshal(categories)
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