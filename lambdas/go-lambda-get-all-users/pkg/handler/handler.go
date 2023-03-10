package handler

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-get-all-users/internal/processor"
	"log"
	"net/http"
)

type GetAllUsersHandler struct {
	processor processor.Processor
}

func NewGetAllUsersHandler(p processor.Processor) *GetAllUsersHandler {
	return &GetAllUsersHandler{processor: p}
}

func (h *GetAllUsersHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("Received request: ", request)
	log.Println("Request QueryStringParameters: ", request.QueryStringParameters)
	filter_name := request.QueryStringParameters["name"]
	filter_surname := request.QueryStringParameters["surname"]
	log.Println("Filter name: ", filter_name)
	log.Println("Filter surname: ", filter_surname)
	users, err := h.processor.GetAllUsers(context.Background(), filter_name)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	response := map[string]interface{}{
		"users":  users.Users,
		"filter": filter_name,
	}

	body, err := json.Marshal(response)
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
