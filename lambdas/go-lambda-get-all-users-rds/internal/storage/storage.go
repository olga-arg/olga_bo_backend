package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go-lambda-get-all-users/pkg/domain"
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

func (r *UserRepository) GetAllUsers(filters map[string]string) (*domain.Users, error) {
	// TODO: Implement filters
	var users *domain.Users
	err := r.db.Scopes(getUserTable()).AutoMigrate(&domain.User{}).Find(&users).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return users, nil
}
