package processor

import (
	"commons/domain"
	"commons/utils/db"
	"context"
	"fmt"
)

type Processor interface {
	UpdateTeam(ctx context.Context, teamID string, newTeam *domain.UpdateTeamRequest, companyId string) error
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

func (p *processor) UpdateTeam(ctx context.Context, teamID string, newTeam *domain.UpdateTeamRequest, companyId string) error {
	if newTeam.AnnualBudget < 0 {
		return fmt.Errorf("annual budget must be greater than 0")
	}
	fmt.Println("Updating team in storage")
	err := p.teamStorage.UpdateTeam(teamID, newTeam, companyId)
	if err != nil {
		fmt.Println("error", err.Error())
		return err
	}
	fmt.Println("Team updated proc")
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
