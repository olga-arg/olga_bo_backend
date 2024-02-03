package processor

import (
	"commons/domain"
	"commons/utils/db"
	"context"
	"go-lambda-get-teams/pkg/dto"
)

type Processor interface {
	GetAllTeams(ctx context.Context, filter map[string]string, companyId string) (*dto.Output, error)
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

func (p *processor) GetAllTeams(ctx context.Context, filter map[string]string, companyId string) (*dto.Output, error) {
	teams, err := p.teamStorage.GetAllTeams(filter, companyId)
	if err != nil {
		return nil, err
	}
	//teams, err = p.teamStorage.GetAllReviewers(teams, companyId)
	return dto.NewOutput(teams), nil
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
