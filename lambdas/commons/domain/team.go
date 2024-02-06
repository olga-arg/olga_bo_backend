package domain

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Team struct {
	ID              string             `json:"id"`
	Name            string             `json:"name"`
	Users           Users              `json:"users" gorm:"many2many:user_teams;"`
	Reviewer        User               `json:"reviewer"`
	ReviewerId      string             `json:"reviewer_id"`
	MonthlySpending int                `json:"monthly_spending"`
	AnnualBudget    int                `json:"annual_budget"`
	Status          ConfirmationStatus `json:"status" default:"Pending"`
	CreatedDate     time.Time          `json:"created_date"`
}

type TeamOutput struct {
	ID              string             `json:"id"`
	Name            string             `json:"name"`
	Users           Users              `json:"users" gorm:"many2many:user_teams;"`
	Reviewer        *User              `json:"reviewer"`
	ReviewerId      string             `json:"reviewer_id"`
	MonthlySpending int                `json:"monthly_spending"`
	AnnualBudget    int                `json:"annual_budget"`
	Status          ConfirmationStatus `json:"status" default:"Pending"`
	CreatedDate     time.Time          `json:"created_date"`
}

type Teams []Team

type DbTeam struct {
	ID                  string             `json:"team_id"`
	Name                string             `json:"team_name"`
	Reviewer            User               `json:"reviewer"`
	ReviewerId          string             `json:"reviewer_id"`
	MonthlySpending     int                `json:"monthly_spending"`
	AnnualBudget        int                `json:"annual_budget"`
	Status              ConfirmationStatus `json:"status"`
	CreatedDate         time.Time          `json:"created_date"`
	UserId              string             `json:"user_id"`
	UserName            string             `json:"user_name"`
	UserSurname         string             `json:"user_surname"`
	UserFullName        string             `json:"user_full_name"`
	UserEmail           string             `json:"user_email"`
	UserMonthlySpending float32            `json:"user_monthly_spending"`
	UserStatus          ConfirmationStatus `json:"user_status"`
	UserRole            UserRoles          `json:"user_role"`
}

type UpdateTeamRequest struct {
	Name         string   `json:"name"`
	AnnualBudget int      `json:"annual_budget"`
	ReviewerId   string   `json:"reviewer_id"`
	AddUsers     []string `json:"add_users"`
	RemoveUsers  []string `json:"remove_users"`
}

type DbTeams []DbTeam

func NewTeam(name, reviewer string, budget int) (*Team, error) {
	var team Team
	id, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("error generating uuid: ", err)
		return nil, err
	}
	team.ID = id.String()
	team.Name = name
	team.ReviewerId = reviewer
	team.AnnualBudget = budget
	team.Status = Created
	team.CreatedDate = time.Now()
	return &team, nil
}
