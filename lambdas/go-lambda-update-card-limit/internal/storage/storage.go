package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go-lambda-update-card-limit/pkg/domain"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func getUserTable(userID string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// TODO: Must receive the user to extract company name to specify the table name
		return tx.Table("users")
	}
}

func (r *UserRepository) UpdateUserCardLimit(newUser *domain.User) error {
	// Save the updated user
	query := r.db.Save(newUser)
	if query.Error != nil {
		fmt.Println("Error updating user card limit:", query.Error)
		return errors.Wrap(query.Error, "failed to update user card limit")
	}
	return nil
}

func (r *UserRepository) GetUserByID(userID string) (*domain.User, error) {
	var user domain.User
	query := r.db.Scopes(getUserTable(userID)).Where("id = ?", userID)
	err := query.First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.Wrap(err, "user not found")
		}
		fmt.Println("Error getting user by ID:", err)
		return nil, errors.Wrap(err, "failed to get user by ID")
	}
	return &user, nil
}
