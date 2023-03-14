package storage

import (
	"github.com/jinzhu/gorm"
	"go-lambda-create-team/pkg/domain"
	"log"
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
		log.Println("Error saving team: ", err)
		return err
	}
	return nil
}
