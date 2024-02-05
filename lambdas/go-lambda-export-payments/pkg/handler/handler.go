package handler

import (
	"commons/utils"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-export-payments/internal/processor"
	"go-lambda-export-payments/pkg/dto"
	"net/http"
)

type ExportPaymentHandler struct {
	processor processor.Processor
}

func NewExportPaymentHandler(p processor.Processor) *ExportPaymentHandler {
	return &ExportPaymentHandler{processor: p}
}

func (h *ExportPaymentHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	var input dto.ExportPaymentsInput
	err = json.Unmarshal([]byte(request.Body), &input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	if len(input.Payments) == 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "No payments to export",
		}, nil
	}

	// Export to csv
	fileUrl, err := h.processor.ExportPayments(companyId, input.Payments)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	// Return the file url
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       fileUrl,
	}, nil
}
