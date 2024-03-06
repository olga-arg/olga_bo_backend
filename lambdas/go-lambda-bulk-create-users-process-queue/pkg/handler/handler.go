package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-bulk-create-users-process-queue/internal/processor"
	"go-lambda-bulk-create-users-process-queue/internal/services"
	"go-lambda-bulk-create-users-process-queue/pkg/dto"
)

type CreateUserServiceHandler struct {
	processor processor.Processor
}

func NewCreateUserHandler(p processor.Processor) *CreateUserServiceHandler {
	return &CreateUserServiceHandler{processor: p}
}

func (h *CreateUserServiceHandler) Handle(context context.Context, event events.SQSEvent) error {
	fmt.Println("Processing queue message")
	fmt.Println("Event:", event)
	for _, message := range event.Records {
		// Decodificar el cuerpo del mensaje en una lista de mapas
		var messages map[string]interface{}
		if err := json.Unmarshal([]byte(message.Body), &messages); err != nil {
			fmt.Println("Error decoding message body:", err)
			continue // Continuar con el próximo mensaje
		}

		// Obtener el company_id del mensaje
		var companyId string
		if len(messages) > 0 {
			companyId = messages["company_id"].(string)
		} else {
			fmt.Println("No company_id found in message")
			return nil
		}

		fmt.Println("Processing company_id:", companyId)
		fmt.Println("Processing message:", messages)

		// Procesar los mensajes de usuario
		user := dto.CreateUserInput{
			Name:    messages["name"].(string),
			Surname: messages["surname"].(string),
			Email:   messages["email"].(string),
			Role:    messages["role"].(string),
		}
		// Llamar a la función createUserService con la información del usuario y el company_id
		if err := h.createUserService(&user, companyId); err != nil {
			fmt.Println("Error processing queue message:", err)
			continue // Continuar con el próximo mensaje
		}
	}
	return nil
}

func (h *CreateUserServiceHandler) createUserService(validUser *dto.CreateUserInput, companyId string) (err error) {
	err = h.processor.CreateUser(context.Background(), validUser, companyId)
	if err != nil {
		fmt.Println("Error creating user:", err)
		return err
	} else {
		err = services.NewDefaultEmailService().SendEmail(validUser.Email, services.Welcome, []string{validUser.Name}, nil)
		if err != nil {
			fmt.Println("Error sending email:", err)
		} else {
			fmt.Println("Email sent successfully to", validUser.Email)
		}
	}

	return err
}
