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

	return nil
}
