package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"go-lambda-create-team/internal/storage"
	"go-lambda-create-team/pkg/domain"
	"go-lambda-create-team/pkg/dto"
	
)

type Processor interface {
	CreateTeam(ctx context.Context, input *dto.CreateTeamInput) error
	ValidateTeamInput(ctx context.Context, input *dto.CreateTeamInput, request events.APIGatewayProxyRequest) error
}

type processor struct {
	storage storage.TeamRepository
}

func New(s storage.TeamRepository) Processor {
	return &processor{
		storage: s,
	}
}

func (p *processor) CreateTeam(ctx context.Context, input *dto.CreateTeamInput) error {
	// Creates a new team
	team, err := domain.NewTeam(input.TeamName, input.ReviewerId, input.AnnualBudget) //input.Employees)
	if err != nil {
		fmt.Println("Error creating team: ", err)
		return err
	}
	// Saves the team to the database if it doesn't already exist
	if err := p.storage.Save(team); err != nil {
		fmt.Println("Error saving team: ", err)
		return err
	}
	// Returns
	return nil
}

func (p *processor) ValidateTeamInput(ctx context.Context, input *dto.CreateTeamInput, request events.APIGatewayProxyRequest) error {
	fmt.Println("Validating input")
	if err := json.Unmarshal([]byte(request.Body), &input); err != nil {
		return fmt.Errorf("invalid request body: %s", err.Error())
	}
	if input.TeamName == "" {
		return fmt.Errorf("team name is required")
	}
	if input.AnnualBudget < 0 {
		return fmt.Errorf("invalid annual budget")
	}
	if request.Body == "" || len(request.Body) < 1 {
		return fmt.Errorf("missing request body")
	}
	// TODO: Validate that reviewer exists in the user table

	// Validate that the team doesn't already exist
	if err := p.storage.GetTeamByName(input.TeamName); err != nil {
		return err
	}
	return nil
}
