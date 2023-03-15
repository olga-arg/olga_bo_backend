package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go-lambda-update-team/pkg/domain"
	"go-lambda-update-team/pkg/dto"
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

func getUsersTable(userID string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// TODO: Must receive the user to extract company name to specify the table name
		return tx.Table("users")
	}
}
func (r *TeamRepository) UpdateTeamBudget(teamID string, newTeam *dto.UpdateTeamRequest) error {
	var team domain.Team
	fmt.Println("Getting team by ID in db")
	err := r.db.Scopes(getTeamTable(teamID)).Where("id = ?", teamID).First(&team).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("No team found")
		return errors.Wrap(err, "No team with that ID found")
	}
	fmt.Println("Team found", team)
	var users []domain.User
	// TODO: Protect against SQL Injection
	err = r.db.Scopes(getUsersTable(teamID)).Where("id IN (?)", newTeam.Users).Find(&users).Error
	fmt.Println("User found", &users)
	if err != nil {
		fmt.Println("Error getting users:", err)
		return errors.Wrap(err, "failed to get users")
	}

	fmt.Println("Team before adding users", team)
	for _, user := range users {
		team.Users = append(team.Users, user)
	}
	fmt.Println("Team after adding users", team)
	err = r.db.Scopes(getTeamTable(teamID)).Save(&team).Error
	if err != nil {
		fmt.Println("Error updating team:", err)
		return errors.Wrap(err, "failed to update team")
	}
	return nil
}
