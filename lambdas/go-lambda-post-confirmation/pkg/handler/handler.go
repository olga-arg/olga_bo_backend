package handler

import (
	"commons/domain"
	"fmt"
	"go-lambda-get-all-users/internal/processor"
)

type PostConfirmationHandler struct {
	processor processor.Processor
}

func NewPostConfirmationHandler(p processor.Processor) *PostConfirmationHandler {
	return &PostConfirmationHandler{processor: p}
}

func (h *PostConfirmationHandler) Handle(request domain.UpdateUserRequest) error {
	// Llama a UpdateUserStatus pasando el company id y el email
	fmt.Println("Request: ", request)
	fmt.Println("Username: ", request.Username)
	err := h.processor.PostConfirmation(request.CompanyId, request.Email, request.Username)
	if err != nil {
		return err
	}

	return nil
}
