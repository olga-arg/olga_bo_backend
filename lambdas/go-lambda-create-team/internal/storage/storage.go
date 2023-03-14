package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"go-lambda-create-team/pkg/domain"
)

type TeamRepository struct {
	Db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{
		Db: db,
	}
}

func getTeamTable(team *domain.Team) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// TODO: Extract company name to specify the table name
		return tx.Table("teams")
	}
}

func (r *TeamRepository) Save(team *domain.Team) error {
	err := r.Db.Scopes(getTeamTable(team)).AutoMigrate(&domain.Team{}).Create(team).Error
	if err != nil {
		fmt.Println("Error saving team: ", err)
		return err
	}
	return nil
}

func (r *TeamRepository) GetTeamByName(teamName string) error {
	var team domain.Team
	query := r.Db.Scopes(getTeamTable(&team)).Where("team_name = ?", teamName)
	err := query.First(&team).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		fmt.Println("Error getting team by name: ", err)
		return err
	}
	return fmt.Errorf("team already exists: %s", teamName)
}
