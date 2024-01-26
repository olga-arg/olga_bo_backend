package db

import (
	"commons/domain"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
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
	fmt.Println("Table name:", companyId)
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

func (r *TeamRepository) DeleteTeam(teamID, companyId string) error {
	var team domain.Team
	fmt.Println("Getting team by ID in db")
	fmt.Println("Team ID:", teamID)
	err := r.Db.Scopes(getTeamTable(companyId)).Where("id = ?", teamID).First(&team).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("No team found")
		return errors.Wrap(err, "No team with that ID found")
	}
	fmt.Println("Team found", team)
	// Validate that team isn't already deleted
	if team.Status == 1 {
		fmt.Println("Team is already deleted")
		return fmt.Errorf("team is already deleted")
	}
	fmt.Println("Team is not deleted")
	// change the team status to deleted
	team.Status = 1
	// Save the updated team
	query := r.Db.Scopes(getTeamTable(companyId)).Save(&team)
	if query.Error != nil {
		fmt.Println("Error deleting team:", query.Error)
		return errors.Wrap(query.Error, "failed to delete team")
	}
	fmt.Println("Team deleted")
	return nil
}

func (r *TeamRepository) Save(team *domain.Team, companyId string) error {
	err := r.Db.Scopes(getTeamTable(companyId)).Create(team).Error
	if err != nil {
		fmt.Println("Error saving team: ", err)
		return err
	}
	return nil
}

//func (r *TeamRepository) FindTeamByID(id string) (*domain.Team, error) {
//	var team domain.Team
//	err := r.Db.Scopes(getTeamsTable(&team)).Where("id = ?", id).First(&team).Error
//	if err != nil {
//		fmt.Println("Error finding team: ", err)
//		return nil, err
//	}
//	return &team, nil
//}

//func (r *TeamRepository) UpdateTeamMonthlySpending(team *domain.Team, paymentAmount float32) error {
//	// Get the current monthly spending of the team
//	var currentMonthlySpending float32
//	currentMonthlySpending = team.MonthlySpending
//
//	// Add the new payment amount to the current monthly spending
//	var newMonthlySpending float32
//	newMonthlySpending = currentMonthlySpending + paymentAmount
//
//	// Save the new monthly spending to the team
//	err := r.Db.Scopes(getTeamsTable(team)).Model(&team).Update("monthly_spending", newMonthlySpending).Error
//	if err != nil {
//		fmt.Println("Error updating team: ", err)
//		return err
//	}
//	return nil
//}
