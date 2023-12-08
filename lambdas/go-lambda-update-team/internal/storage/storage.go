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
	db.AutoMigrate(&domain.Team{})
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
	err := r.db.Scopes(getTeamTable(teamID)).Where("id = ?", teamID).Preload("Users").First(&team).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("No team found")
		return errors.Wrap(err, "No team with that ID found")
	}
	var reviewer domain.User
	if newTeam.ReviewerId != "" && newTeam.ReviewerId != team.ReviewerId {
		println("Updating reviewer")
		// validate that newTeam.ReviewerId is a valid user
		err = r.db.Scopes(getUsersTable(teamID)).Where("id = ?", newTeam.ReviewerId).First(&reviewer).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("No reviewer found")
		}
		team.ReviewerId = newTeam.ReviewerId
	}

	var users []domain.User
	// TODO: Protect against SQL Injection
	if len(newTeam.AddUsers) > 0 {
		err = r.db.Scopes(getUsersTable(teamID)).Where("id IN (?)", newTeam.AddUsers).Find(&users).Error
		if err != nil {
			fmt.Println("Error getting users:", err)
			return errors.Wrap(err, "failed to get users")
		}
		for _, user := range users {
			team.Users = append(team.Users, user)
		}
		// Actualizar la relación en la tabla intermedia
		err = r.db.Model(&team).Association("Users").Append(users).Error
		if err != nil {
			fmt.Println("Error adding users to team:", err)
			return errors.Wrap(err, "failed to add users to team")
		}
	}
	if len(newTeam.RemoveUsers) > 0 {
		var remainingUsers []domain.User
		for _, teamUser := range team.Users {
			shouldKeep := true
			for _, removeUserID := range newTeam.RemoveUsers {
				if teamUser.ID == removeUserID {
					shouldKeep = false
					break
				}
			}
			if shouldKeep {
				remainingUsers = append(remainingUsers, teamUser)
			}
		}
		if reviewer.Email != "" {
			remainingUsers = append(remainingUsers, reviewer)
		} else {
			// search in db the user with the reviewer id
			err = r.db.Scopes(getUsersTable(teamID)).Where("id = ?", team.ReviewerId).First(&reviewer).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Println("No reviewer found")
			}
			remainingUsers = append(remainingUsers, reviewer)
		}
		team.Users = remainingUsers
		// Eliminar las relaciones en la tabla intermedia
		err = r.db.Model(&team).Association("Users").Replace(remainingUsers).Error
		if err != nil {
			fmt.Println("Error removing users from team:", err)
			return errors.Wrap(err, "failed to remove users from team")
		}
	}

	// Actualizar el equipo en la base de datos
	err = r.db.Scopes(getTeamTable(teamID)).Save(&team).Error
	if err != nil {
		fmt.Println("Error updating team:", err)
		return errors.Wrap(err, "failed to update team")
	}
	return nil
}
