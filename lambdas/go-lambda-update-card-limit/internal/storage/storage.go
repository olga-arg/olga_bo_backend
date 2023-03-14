package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go-lambda-update-card-limit/pkg/domain"
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

func getUserTable(userID string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// TODO: Must receive the user to extract company name to specify the table name
		return tx.Table("users")
	}
}

func (r *UserRepository) UpdateUserCardLimit(userID string, purchaseLimit int, monthlyLimit int) (*domain.User, error) {
	// Get the current user
	user := &domain.User{}
	err := r.db.Scopes(getUserTable(userID)).Where("id = ?", userID).First(user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.Wrap(err, "user not found")
		}
		log.Println("Error getting user by ID:", err)
		return nil, errors.Wrap(err, "failed to get user by ID")
	}

	// Update only the specified field
	if purchaseLimit != -1 {
		user.PurchaseLimit = purchaseLimit
	}
	if monthlyLimit != -1 {
		user.MonthlyLimit = monthlyLimit
	}

	// Save the updated user
	result := r.db.Save(user)
	if result.Error != nil {
		log.Println("Error updating user card limit:", result.Error)
		return nil, errors.Wrap(result.Error, "failed to update user card limit")
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(userID string) (*domain.User, error) {
	var user domain.User
	query := r.db.Scopes(getUserTable(userID)).Where("id = ?", userID)
	err := query.First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.Wrap(err, "user not found")
		}
		log.Println("Error getting user by ID:", err)
		return nil, errors.Wrap(err, "failed to get user by ID")
	}
	return &user, nil
}
