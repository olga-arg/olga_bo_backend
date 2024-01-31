package handler

import (
	"commons/utils"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-create-category/internal/processor"
	"go-lambda-create-category/pkg/dto"
	"net/http"
)

type CreateCategoryHandler struct {
	processor processor.Processor
}

func NewCreateCategoryHandler(p processor.Processor) *CreateCategoryHandler {
	return &CreateCategoryHandler{processor: p}
}
func (h *CreateCategoryHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input dto.CreateCategoryInput

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

	// Validate input
	err = h.processor.ValidateCategoryInput(context.Background(), &input, request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	err = h.processor.CreateCategory(context.Background(), &input, companyId)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
		},
		Body: "Category created successfully",
	}, nil
}
