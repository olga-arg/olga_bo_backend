package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go-lambda-update-team/pkg/domain"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func getTeamTable(teamID string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// TODO: Must receive the user to extract company name to specify the table name
		return tx.Table("teams")
	}
}
func (r *TeamRepository) UpdateTeamBudget(teamID string, annualBudget int) error {
	var team domain.Team
	fmt.Println("Getting team by ID in db")
	err := r.db.Scopes(getTeamTable(teamID)).Where("id = ?", teamID).First(&team).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("No team found")
		return errors.Wrap(err, "No team with that ID found")
	}
	fmt.Println("Team found", team)

	// change the team annual budget
	team.AnnualBudget = annualBudget
	// Save the updated team
	query := r.db.Save(team)
	if query.Error != nil {
		fmt.Println("Error updating team:", query.Error)
		return errors.Wrap(query.Error, "failed to update team")
	}
	fmt.Println("Team updated")
	return nil
}
