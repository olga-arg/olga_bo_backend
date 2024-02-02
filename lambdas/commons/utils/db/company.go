package db

import (
	"commons/domain"
	"fmt"
	"github.com/jinzhu/gorm"
)

type CompanyRepository struct {
	Db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) *CompanyRepository {
	db.AutoMigrate(&domain.Company{})
	return &CompanyRepository{
		Db: db,
	}
}

func (r *CompanyRepository) Save(company *domain.Company) error {
	db := r.Db.AutoMigrate(&domain.Company{})
	err := db.Save(company).Error
	if err != nil {
		fmt.Println("Error saving company: ", err)
		return err
	}
	return nil
}

func (r *CompanyRepository) CreateCompanySpecificTables(companyId string) error {

	// Crear tabla de usuarios para la empresa
	err := r.Db.Table(fmt.Sprintf("%s_users", companyId)).AutoMigrate(&domain.User{}).Error
	if err != nil {
		fmt.Println("Error creating user table: ", err)
		return err
	}
	// Crear tabla de pagos para la empresa
	err = r.Db.Table(fmt.Sprintf("%s_payments", companyId)).AutoMigrate(&domain.Payment{}).Error
	if err != nil {
		fmt.Println("Error creating payment table: ", err)
		return err
	}

	// Crear tabla de equipos para la empresa
	err = r.Db.Table(fmt.Sprintf("%s_teams", companyId)).AutoMigrate(&domain.Team{}).Error
	if err != nil {
		fmt.Println("Error creating team table: ", err)
		return err
	}

	// Crear la tabla users_teams para la empresa
	usersTeamsTableName := fmt.Sprintf("%s_users_teams", companyId)
	usersTableName := fmt.Sprintf("%s_users", companyId)
	teamsTableName := fmt.Sprintf("%s_teams", companyId)

	// Migrar la tabla users_teams
	err = r.Db.Table(usersTeamsTableName).AutoMigrate(&domain.UserTeam{}).Error
	if err != nil {
		fmt.Println("Error creating users_teams table: ", err)
		return err
	}

	// Añadir foreign keys después de crear la tabla
	// Añadir foreign keys después de crear la tabla
	err = r.Db.Debug().Exec(fmt.Sprintf(`
	ALTER TABLE %s
	ADD CONSTRAINT fk_users
	FOREIGN KEY (user_id) REFERENCES %s(id);
`, r.Db.Dialect().Quote(usersTeamsTableName), r.Db.Dialect().Quote(usersTableName))).Error

	if err != nil {
		fmt.Println("Error adding foreign key for users: ", err)
		return err
	}

	err = r.Db.Debug().Exec(fmt.Sprintf(`
	ALTER TABLE %s
	ADD CONSTRAINT fk_teams
	FOREIGN KEY (team_id) REFERENCES %s(id);
`, r.Db.Dialect().Quote(usersTeamsTableName), r.Db.Dialect().Quote(teamsTableName))).Error

	if err != nil {
		fmt.Println("Error adding foreign key for teams: ", err)
		return err
	}

	// Add the unique key constraint for the columns user_id and team_id
	err = r.Db.Debug().Exec(fmt.Sprintf(`
	ALTER TABLE %s
	ADD CONSTRAINT unique_user_team UNIQUE (user_id, team_id);
`, r.Db.Dialect().Quote(usersTeamsTableName))).Error

	if err != nil {
		fmt.Println("Error adding unique key for users_teams: ", err)
		return err
	}

	return nil
}
