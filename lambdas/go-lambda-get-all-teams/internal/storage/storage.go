package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go-lambda-get-all-teams/pkg/domain"
	"log"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
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
	query := r.db.Scopes(getTeamTable())

	// TODO: Always filter by confirmed teams
	// Apply filters to the query
	if teamName, ok := filters["team_name"]; ok {
		query = query.Where("team_name ILIKE ?", "%"+teamName+"%")
	}

	if annualBudget, ok := filters["annual_budget"]; ok {
		query = query.Where("annual_budget = ?", annualBudget)
	}

	// Execute the query
	err := query.Find(&teams).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("No teams found")
		return nil, nil
	}
	if err != nil {
		log.Println("Error getting teams:", err)
		return nil, err
	}

	return teams, nil
}