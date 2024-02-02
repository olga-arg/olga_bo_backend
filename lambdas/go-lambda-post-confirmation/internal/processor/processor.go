package processor

import (
	"commons/utils/db"
)

type Processor interface {
	PostConfirmation(companyId, email, username string) error
}

type processor struct {
	userStorage *db.UserRepository
}

func NewProcessor(storage *db.UserRepository) Processor {
	return &processor{
		userStorage: storage,
	}
}

func (p *processor) PostConfirmation(companyId, email, username string) error {
	err := p.userStorage.UpdateUserStatus(companyId, email)
	if err != nil {
		return err
	}

	// Update user email to verified in cognito
	err = p.userStorage.UpdateEmailVerified(username)

	return nil
}
