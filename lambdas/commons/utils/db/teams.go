package db

import (
	"commons/domain"
	"fmt"
	"github.com/jinzhu/gorm"
)

type TeamRepository struct {
	Db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	db.AutoMigrate(&domain.Team{})
	return &TeamRepository{
		Db: db,
	}
}

func getTeamTable(companyId string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tableName := fmt.Sprintf("%s_teams", companyId)
		return tx.Table(tableName)
	}
}

func (r *TeamRepository) GetTeamByName(name, companyId string) error {
	// it should only return an error if the team already exists if it doesn't exist, it should return nil
	var team domain.Team
	err := r.Db.Scopes(getTeamTable(companyId)).Where("name = ?", name).First(&team).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		fmt.Println("Error getting team by name: ", err)
		return err
	}
	return fmt.Errorf("team already exists: %s", name)
}

func (r *TeamRepository) Save(team *domain.Team, companyId string) error {
	err := r.Db.Scopes(getTeamTable(companyId)).Create(team).Error
	if err != nil {
		fmt.Println("Error saving team: ", err)
		return err
	}
	return nil
}
