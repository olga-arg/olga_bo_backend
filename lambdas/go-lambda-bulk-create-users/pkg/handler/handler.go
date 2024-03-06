package handler

import (
	"commons/domain"
	"commons/utils"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-bulk-create-users/internal/processor"
	"go-lambda-bulk-create-users/internal/services"
	"go-lambda-bulk-create-users/pkg/dto"
	"net/http"
)

type CreateUserHandler struct {
	processor processor.Processor
}

func NewCreateUserHandler(p processor.Processor) *CreateUserHandler {
	return &CreateUserHandler{processor: p}
}

func (h *CreateUserHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	// Validate user
	allowedRoles := []domain.UserRoles{domain.Admin}
	isAuthorized, err := h.processor.ValidateUser(context.Background(), email, companyId, allowedRoles)
	if err != nil {
		fmt.Println("Error validating user: ", err)
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

	users, err := h.processor.ParseCSVFromRequest(context.Background(), request)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: err.Error()}, nil
	}

	var failedUsers []domain.UserNotCreated
	var validUsers []dto.CreateUserInput

	for _, userInput := range users {
		fmt.Println("Validating user input: ", userInput)
		err = h.processor.ValidateUserInput(context.Background(), &userInput)
		if err != nil {
			failedUsers = append(failedUsers, domain.UserNotCreated{Email: userInput.Email, Reason: err.Error()})
			continue
		}
		// Agrega el usuario válido a la lista de usuarios válidos
		validUsers = append(validUsers, userInput)
	}

	fmt.Println("Creating users:", validUsers)

	// Crear un llamado a la cola de SQS por cada usuario válido

	err = services.SendMessageToQueue(validUsers, companyId)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "Users sent to queue successfully",
	}, nil

}
