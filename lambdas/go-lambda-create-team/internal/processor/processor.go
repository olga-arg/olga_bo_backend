package processor

import (
	"commons/domain"
	"commons/utils/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-create-team/pkg/dto"
)

type Processor interface {
	CreateTeam(ctx context.Context, input *dto.CreateTeamInput, companyId string) error
	ValidateTeamInput(ctx context.Context, input *dto.CreateTeamInput, request events.APIGatewayProxyRequest, companyId string) error
	ValidateUser(ctx context.Context, email, companyId string, allowedRoles []domain.UserRoles) (bool, error)
}

type processor struct {
	teamStorage db.TeamRepository
	userStorage db.UserRepository
}

func New(s db.TeamRepository, u db.UserRepository) Processor {
	return &processor{
		teamStorage: s,
		userStorage: u,
	}
}

func (p *processor) CreateTeam(ctx context.Context, input *dto.CreateTeamInput, companyId string) error {
	// Creates a new team
	team, err := domain.NewTeam(input.Name, input.ReviewerId, input.AnnualBudget) //input.Employees)
	if err != nil {
		fmt.Println("Error creating team: ", err)
		return err
	}
	// Saves the team to the database if it doesn't already exist
	if err := p.teamStorage.Save(team, companyId); err != nil {
		fmt.Println("Error saving team: ", err)
		return err
	}
	// Returns
	return nil
}

func (p *processor) ValidateTeamInput(ctx context.Context, input *dto.CreateTeamInput, request events.APIGatewayProxyRequest, companyId string) error {
	fmt.Println("Validating input")
	if err := json.Unmarshal([]byte(request.Body), &input); err != nil {
		return fmt.Errorf("invalid request body: %s", err.Error())
	}
	if input.Name == "" {
		return fmt.Errorf("team name is required")
	}
	if input.AnnualBudget < 0 {
		return fmt.Errorf("invalid annual budget")
	}
	if request.Body == "" || len(request.Body) < 1 {
		return fmt.Errorf("missing request body")
	}
	// Validate that the team doesn't already exist
	if err := p.teamStorage.GetTeamByName(input.Name, companyId); err != nil {
		return err
	}
	// Validate that the reviewer exists only if provided
	if input.ReviewerId != "" {
		err := p.teamStorage.GetReviewerById(input.ReviewerId, companyId)
		if err != nil {
			return err
		}
	}
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
