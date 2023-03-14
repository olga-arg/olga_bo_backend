package handler

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-create-user/internal/processor"
	"go-lambda-create-user/internal/services"
	"go-lambda-create-user/pkg/dto"

	"net/http"
)

type CreateUserHandler struct {
	processor processor.Processor
}

func NewCreateUserHandler(p processor.Processor) *CreateUserHandler {
	return &CreateUserHandler{processor: p}
}

func (h *CreateUserHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var input dto.CreateUserInput

	// Validate input
	fmt.Println("Validating input")
	err := h.processor.ValidateUserInput(context.Background(), &input, request)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	err = h.processor.CreateUser(context.Background(), &input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	body := "Hola " + input.Name + ",\n\n¡Bienvenido a Olga!\nAhora tienes acceso a la mejor forma de administrar los gastos en tu empresa.\n\nPara empezar, descarga nuestra app móvil en tu dispositivo:\n\n[enlace para descargar la app iOS/Android]\n\nUna vez que hayas descargado la aplicación, haz click en el siguiente enlace para comenzar tu procesa de registro y comenzar a disfrutar de los siguientes beneficios:\nAprobación de gastos de manera instantanea\nTarjetas Fisicas o Virtuales para ti y tus compañeros de trabajo\nNada de guardar el ticket para presentarlo en contabilidad, ahora solo basta una foto en el momento de la compra y listo!\nNo dudes en ponerte en contacto con nuestro equipo de soporte si tienes alguna pregunta o necesitas ayuda para empezar.\n\nSaludos cordiales,\nEl equipo de Olga"
	fmt.Println("Sending email from handler to: " + input.Email)
	services.NewDefaultEmailService().SendEmail("Bienvenido a Olga :)", body, input.Email, nil)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "User created successfully",
	}, nil
}
