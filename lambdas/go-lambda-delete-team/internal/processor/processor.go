package processor

import (
	"commons/domain"
	"commons/utils/db"
	"context"
	"fmt"
)

type Processor interface {
	DeleteTeam(ctx context.Context, teamID, companyId string) error
	ValidateUser(ctx context.Context, email, companyId string, allowedRoles []domain.UserRoles) (bool, error)
}

type processor struct {
	teamStorage db.TeamRepository
	userStorage db.UserRepository
}

func NewProcessor(teamStorage db.TeamRepository, userStorage db.UserRepository) Processor {
	return &processor{
		teamStorage: teamStorage,
		userStorage: userStorage,
	}
}

func (p *processor) DeleteTeam(ctx context.Context, teamID, companyId string) error {
	fmt.Println("Deleting team in storage")
	err := p.teamStorage.DeleteTeam(teamID, companyId)
	if err != nil {
		fmt.Println("error", err.Error())
		return err
	}
	fmt.Println("Team deleted proc")
	return nil
}

func (p *processor) ValidateUser(ctx context.Context, email, companyId string, allowedRoles []domain.UserRoles) (bool, error) {
	// Validate user
	isAuthorized, err := p.userStorage.IsUserAuthorized(email, companyId, allowedRoles)
	if err != nil {
		return false, err
	}
	if isAuthorized {
		return true, nil
	}
	return false, nil
}
