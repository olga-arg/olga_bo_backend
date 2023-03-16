package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go-lambda-get-teams/pkg/domain"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	db.AutoMigrate(&domain.Team{})
	return &TeamRepository{
		db: db,
	}
}
func getTeamTable() func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// TODO: Must receive the user to extract company name to specify the table name
		return tx.Table("teams")
	}
}

func (r *TeamRepository) GetAllTeams(filters map[string]string) ([]domain.Team, error) {
	var teams []domain.Team
	query := r.db.Scopes(getTeamTable()).Preload("Users").Where("status = ?", 0)
	// Apply filters to the query
	if teamName, ok := filters["team_name"]; ok {
		query = query.Where("team_name ILIKE ?", "%"+teamName+"%")
	}

	if annualBudget, ok := filters["annual_budget"]; ok {
		query = query.Where("annual_budget = ?", annualBudget)
	}

	// Execute the query
	err := query.Find(&teams).Error
	fmt.Println("teams", teams)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("No teams found")
		return nil, nil
	}
	if err != nil {
		fmt.Println("Error getting teams:", err)
		return nil, err
	}

	return teams, nil
}
