package processor

import (
	"commons/domain"
	"commons/utils/db"
	"encoding/json"
)

type Processor interface {
	PostConfirmation(requestBody string) error
}

type processor struct {
	userStorage *db.UserRepository
}

func NewProcessor(storage *db.UserRepository) Processor {
	return &processor{
		userStorage: storage,
	}
}

func (p *processor) PostConfirmation(request string) error {
	// Define una estructura para el JSON, por ejemplo:
	var input domain.UpdateUserRequest

	// Deserializa el JSON en la estructura
	err := json.Unmarshal([]byte(request), &input)
	if err != nil {
		return err
	}

	// Llama a UpdateUserStatus pasando el company id y el email
	err = p.userStorage.UpdateUserStatus(input.CompanyId, input.Email)
	if err != nil {
		return err
	}

	return nil
}
