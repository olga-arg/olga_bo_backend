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
	db.AutoMigrate(&domain.Team{})
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

func getUserTable(user *domain.User) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Table("users")
	}
}

func (r *TeamRepository) Save(team *domain.Team) error {
	err := r.Db.Scopes(getTeamTable(team)).Create(team).Error
	if err != nil {
		fmt.Println("Error saving team: ", err)
		return err
	}
	return nil
}

func (r *TeamRepository) GetTeamByName(name string) error {
	// it should only return an error if the team already exists if it doesn't exist, it should return nil
	var team domain.Team
	err := r.Db.Scopes(getTeamTable(&team)).Where("name = ?", name).First(&team).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
		fmt.Println("Error getting team by name: ", err)
		return err
	}
	return fmt.Errorf("team already exists: %s", name)
}

func (r *TeamRepository) GetReviewerById(id string) error {
	var user domain.User
	err := r.Db.Scopes(getUserTable(&user)).Where("id = ?", id).First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return fmt.Errorf("user not found: %s", id)
		}
		fmt.Println("Error getting user by id: ", err)
		return err
	}
	return nil
}
