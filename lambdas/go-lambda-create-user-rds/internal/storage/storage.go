package storage

import (
	"github.com/jinzhu/gorm"
	"go-lambda-create-user/pkg/domain"
	"log"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func getUserTable(user *domain.User) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// Extract company name to specify the table name
		return tx.Table("users")
	}
}

func (r *UserRepository) EmailAlreadyExists(email string) (bool, error) {
	result := r.db.Scopes(getUserTable(&domain.User{})).Where("email = ?", email).Find(&domain.User{})
	if result.Error != nil {
		log.Println("Error checking if email already exists: ", result.Error)
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *UserRepository) Save(user *domain.User) error {
	defer r.db.Close()
	r.db.Scopes(getUserTable(user)).AutoMigrate(&domain.User{}).Create(user)
	return nil
}
