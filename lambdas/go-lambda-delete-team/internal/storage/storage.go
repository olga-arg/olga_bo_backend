package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go-lambda-delete-team/pkg/domain"
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

func getTeamTable(teamID string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// TODO: Must receive the user to extract company name to specify the table name
		return tx.Table("teams")
	}
}

func (r *TeamRepository) DeleteTeam(newTeam *domain.Team) error {
	// change the team status to deleted
	newTeam.Status = 1
	// Save the updated team
	query := r.db.Save(newTeam)
	if query.Error != nil {
		log.Println("Error deleting team:", query.Error)
		return errors.Wrap(query.Error, "failed to delete team")
	}
	return nil
}

func (r *TeamRepository) GetTeamByID(teamID string) (*domain.Team, error) {
	var team domain.Team
	query := r.db.Scopes(getTeamTable(teamID)).Where("id = ?", teamID)
	err := query.First(&team).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.Wrap(err, "team not found")
		}
		log.Println("Error getting team by ID:", err)
		return nil, errors.Wrap(err, "failed to get team by ID")
	}
	return &team, nil
}
