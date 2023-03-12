package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go-lambda-get-all-users-rds/pkg/domain"
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
func getUserTable() func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// TODO: Must receive the user to extract company name to specify the table name
		return tx.Table("users")
	}
}

func (r *UserRepository) GetAllUsers(filters map[string]string) ([]domain.User, error) {
	// TODO: Implement filters
	var users []domain.User
	log.Println("Getting all users")
	err := r.db.Scopes(getUserTable()).Find(&users).Error
	log.Println("Users found: ", users)
	log.Println("Error: ", err)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("No users found")
		return nil, nil
	}
	log.Println("returning users")
	return users, nil
}
