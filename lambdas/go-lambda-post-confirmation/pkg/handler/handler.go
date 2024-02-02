package handler

import (
	"commons/domain"
	"encoding/json"
	"fmt"
	"go-lambda-get-all-users/internal/processor"
)

type PostConfirmationHandler struct {
	processor processor.Processor
}

func NewPostConfirmationHandler(p processor.Processor) *PostConfirmationHandler {
	return &PostConfirmationHandler{processor: p}
}

func (h *PostConfirmationHandler) Handle(request string) error {
	fmt.Println("HOLA LA PUTA MADRE COMO ESTAS AMIGO?")
	// Define una estructura para el JSON, por ejemplo:
	var input domain.UpdateUserRequest

	fmt.Println("Request: ", request)

	// Deserializa el JSON en la estructura
	err := json.Unmarshal([]byte(request), &input)
	if err != nil {
		fmt.Println("Error unmarshalling request: ", request)
		return err
	}

	// Llama a UpdateUserStatus pasando el company id y el email
	err = h.processor.PostConfirmation(input.CompanyId)
	if err != nil {
		return err
	}

	return nil
}
