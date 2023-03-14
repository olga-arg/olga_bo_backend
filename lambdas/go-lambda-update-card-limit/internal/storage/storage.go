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
	user := r.db.Scopes(getUserTable(userID)).Where("id = ?", userID)

	data := make(map[string]interface{})
	if purchaseLimit != -1 {
		data["purchase_limit"] = purchaseLimit
	}
	if monthlyLimit != -1 {
		data["monthly_limit"] = monthlyLimit
	}

	var updatedUser domain.User
	result := user.Updates(data)
	if result.Error != nil {
		log.Println("Error updating user card limit:", result.Error)
		return nil, errors.Wrap(result.Error, "failed to update user card limit")
	}
	if result.RowsAffected == 0 {
		log.Println("No user card found for update")
		return nil, errors.New("no user card found for update")
	}
	err := user.First(&updatedUser).Error
	if err != nil {
		log.Println("Error getting user after update:", err)
		return nil, errors.Wrap(err, "failed to get user after update")
	}
	return &updatedUser, nil
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
