package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go-lambda-create-user/pkg/domain"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	db.AutoMigrate(&domain.User{})
	return &UserRepository{
		db: db,
	}
}

func getUserTable(user *domain.User) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// TODO: Extract company name to specify the table name
		return tx.Table("users")
	}
}

func (r *UserRepository) EmailAlreadyExists(email string) (bool, error) {
	err := r.db.Scopes(getUserTable(&domain.User{})).Preload("Teams").Where("email = ?", email).First(&domain.User{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return true, nil
}

func (r *UserRepository) Save(user *domain.User) error {
	err := r.db.Scopes(getUserTable(user)).Create(user).Error
	if err != nil {
		fmt.Println("Error saving user: ", err)
		return err
	}
	return nil
}
