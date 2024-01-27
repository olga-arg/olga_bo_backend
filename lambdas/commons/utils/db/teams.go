package db

import (
	"commons/domain"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"strconv"
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

func getUserTeamTable(companyId string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tableName := fmt.Sprintf("%s_users_teams", companyId)
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

func (r *TeamRepository) GetAllTeams(filters map[string]string, companyId string) (domain.DbTeams, error) {
	var teams domain.DbTeams

	usersTeamsTableName := fmt.Sprintf("%s_users_teams", companyId)
	usersTableName := fmt.Sprintf("%s_users", companyId)
	teamsTableName := fmt.Sprintf("%s_teams", companyId)
	querySyntax := fmt.Sprintf(
		`SELECT teams.id, teams.name, teams.monthly_spending as "team_monthly_spending", teams.annual_budget as "team_annual_budget", teams.status as "team_status", teams.created_date as "team_created_date", 
	users.id as "user_id", users.name as "user_name", users.surname as "user_surname", users.email as "user_email", users.monthly_spending as "user_monthly_spending"
    FROM "%s" as teams 
    JOIN "%s" as users_teams ON teams.id = users_teams.team_id 
    JOIN "%s" as users ON users.id = users_teams.user_id`,
		teamsTableName, usersTeamsTableName, usersTableName)

	if teamName, ok := filters["name"]; ok {
		println("team name:", teamName)
		querySyntax = querySyntax + fmt.Sprintf(" WHERE teams.name = '%s'", teamName)
	} else if annualBudget, ok := filters["annual_budget"]; ok {
		println("annual budget:", annualBudget)
		Int, err := strconv.Atoi(annualBudget)
		if err != nil {
			fmt.Println("Error converting annual budget to int:", err)
		} else {
			querySyntax = querySyntax + fmt.Sprintf(" WHERE teams.annual_budget = %d", Int)
		}
	}
	query := r.Db.Raw(querySyntax)
	query = query.Order("teams.name ASC")
	// Ejecutar la consulta
	err := query.Scan(&teams).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("No teams found")
		return nil, nil
	}
	if err != nil {
		fmt.Println("Error getting teams:", err)
		return nil, err
	}

	return teams, nil
}

func (r *TeamRepository) GetAllReviewers(teams []domain.Team, companyId string) ([]domain.Team, error) {
	for i, team := range teams {
		var reviewer domain.User
		err := r.Db.Scopes(getUserTable(companyId)).Model(&reviewer).Where("id = ?", team.ReviewerId).Find(&reviewer).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("No reviewer found for team:", team)
		}
		if err != nil {
			fmt.Println("Error getting reviewer:", err)
		}
		teams[i].Reviewer = reviewer
	}
	return teams, nil
}

func (r *TeamRepository) Save(team *domain.Team, companyId string) error {
	err := r.Db.Scopes(getTeamTable(companyId)).Create(team).Error
	if err != nil {
		fmt.Println("Error saving team: ", err)
		return err
	}
	return nil
}

func (r *TeamRepository) UpdateTeamBudget(teamID string, newTeam *domain.UpdateTeamRequest, companyId string) error {
	var team domain.Team
	fmt.Println("Getting team by ID in db")
	err := r.Db.Scopes(getTeamTable(companyId)).Where("id = ?", teamID).
		Preload("Users", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(getUserTable(companyId))
		}).First(&team).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("No team found")
		return errors.Wrap(err, "No team with that ID found")
	}
	var reviewer domain.User
	if newTeam.ReviewerId != "" && newTeam.ReviewerId != team.ReviewerId {
		println("Updating reviewer")
		// validate that newTeam.ReviewerId is a valid user
		err = r.Db.Scopes(getUserTable(companyId)).Where("id = ?", newTeam.ReviewerId).First(&reviewer).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("No reviewer found")
		}
		team.ReviewerId = newTeam.ReviewerId
	}

	var users []domain.User
	// TODO: Protect against SQL Injection
	if len(newTeam.AddUsers) > 0 {
		err = r.Db.Scopes(getUserTable(companyId)).Where("id IN (?)", newTeam.AddUsers).Find(&users).Error
		if err != nil {
			fmt.Println("Error getting users:", err)
			return errors.Wrap(err, "failed to get users")
		}
		for _, user := range users {
			team.Users = append(team.Users, user)
		}
		// Actualizar la relación en la tabla intermedia
		err = r.Db.Model(&team).Association("Users").Append(users).Error
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
			err = r.Db.Scopes(getUserTable(companyId)).Where("id = ?", team.ReviewerId).First(&reviewer).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Println("No reviewer found")
			}
			remainingUsers = append(remainingUsers, reviewer)
		}
		team.Users = remainingUsers
		// Eliminar las relaciones en la tabla intermedia
		err = r.Db.Model(&team).Association("Users").Replace(remainingUsers).Error
		if err != nil {
			fmt.Println("Error removing users from team:", err)
			return errors.Wrap(err, "failed to remove users from team")
		}
	}

	// Actualizar el equipo en la base de datos
	err = r.Db.Scopes(getTeamTable(companyId)).Save(&team).Error
	if err != nil {
		fmt.Println("Error updating team:", err)
		return errors.Wrap(err, "failed to update team")
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
