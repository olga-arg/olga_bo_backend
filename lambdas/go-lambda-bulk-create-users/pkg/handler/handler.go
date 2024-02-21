package handler

import (
	"commons/domain"
	"commons/utils"
	"context"
	"encoding/json"
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

	// Aquí asumimos que tienes una función para parsear el CSV del request y devolver una lista de CreateUserInput
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

	// Ahora, crea todos los usuarios válidos en un solo llamado
	var cognitoFailedUsers []domain.UserNotCreated
	if len(validUsers) > 0 {
		cognitoFailedUsers, err = h.processor.CreateMultipleUsers(context.Background(), validUsers, companyId)
		if err != nil {
			// Manejar error
			fmt.Println("Error creating multiple users:", err)
			// Puedes elegir agregar todos los usuarios válidos al listado de fallos con una razón general
			// o manejar este error de otra manera dependiendo de tu lógica de negocio.
		}
	}

	fmt.Println("Cognito failed users:", cognitoFailedUsers)

	// Agrega los usuarios que fallaron en Cognito a la lista de fallos
	failedUsers = append(failedUsers, cognitoFailedUsers...)

	fmt.Println("Failed users:", failedUsers)

	// Filtra los usuarios válidos para excluir a los que fallaron en Cognito
	var usersForEmail []dto.CreateUserInput
	for _, validUser := range validUsers {
		var failed bool
		for _, failedUser := range cognitoFailedUsers {
			if validUser.Email == failedUser.Email {
				failed = true
				break
			}
		}
		if !failed {
			usersForEmail = append(usersForEmail, validUser)
		}
	}

	fmt.Println("Users for email:", usersForEmail)

	// Envía correos electrónicos solo a los usuarios que no fallaron en la creación de Cognito
	for _, user := range usersForEmail {
		err = services.NewDefaultEmailService().SendEmail(user.Email, services.Welcome, []string{user.Name}, nil)
		if err != nil {
			fmt.Println("Error sending email:", err)
		} else {
			fmt.Println("Email sent successfully to", user.Email)
		}
	}

	result := domain.CreateUserResult{
		FailedUsers:  failedUsers,
		SuccessCount: len(usersForEmail),
	}

	fmt.Println("Result:", result)

	// Serializar result a JSON
	resultBytes, err := json.Marshal(result)
	if err != nil {
		// Manejar error de serialización aquí
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error serializing response",
		}, nil
	}

	fmt.Println("Result bytes:", string(resultBytes))

	// Construir y retornar APIGatewayProxyResponse
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK, // Asumiendo éxito; ajusta según sea necesario
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(resultBytes),
	}, nil
}
