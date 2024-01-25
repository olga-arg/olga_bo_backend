package db

import (
	"commons/domain"
	"fmt"
	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	Db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	db.AutoMigrate(&domain.User{})
	return &UserRepository{
		Db: db,
	}
}

func getUserTable(companyID string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tableName := fmt.Sprintf("%s_users", companyID)
		return tx.Table(tableName)
	}
}

func (r *UserRepository) GetUserIdByEmail(email string, companyID string) (*domain.User, error) {
	var user domain.User
	err := r.Db.Scopes(getUserTable(companyID)).Where("email = ?", email).First(&user).Error
	if err != nil {
		fmt.Println("Error getting user id: ", err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(newUser *domain.User, companyId string) error {
	err := r.Db.Scopes(getUserTable(companyId)).Save(newUser).Error
	if err != nil {
		fmt.Println("Error updating user:", err)
		return err
	}
	return nil
}

func (r *UserRepository) Save(user *domain.User, companyId string) error {
	err := r.Db.Scopes(getUserTable(companyId)).Create(user).Error
	if err != nil {
		fmt.Println("Error saving user: ", err)
		return err
	}
	return nil
}
