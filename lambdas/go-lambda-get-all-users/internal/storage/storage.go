package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go-lambda-get-all-users/pkg/domain"
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
	var users []domain.User
	query := r.db.Scopes(getUserTable())

	// TODO: Always filter by confirmed users
	// Apply filters to the query
	if fullName, ok := filters["name"]; ok {
		query = query.Where("full_name ILIKE ?", "%"+fullName+"%")
	}
	if email, ok := filters["email"]; ok {
		query = query.Where("email ILIKE ?", "%"+email+"%")
	}
	if isAdmin, ok := filters["isAdmin"]; ok {
		query = query.Where("is_admin = ?", isAdmin)
	}

	// Execute the query
	err := query.Find(&users).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("No users found")
		return nil, nil
	}
	if err != nil {
		log.Println("Error getting users:", err)
		return nil, err
	}

	return users, nil
}
