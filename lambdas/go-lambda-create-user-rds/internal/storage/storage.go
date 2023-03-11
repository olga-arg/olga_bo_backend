package storage

import (
	"github.com/jinzhu/gorm"

	"go-lambda-create-user/pkg/domain"
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

func (r *UserRepository) Save(user *domain.User) error {
	defer r.db.Close()
	r.db.Scopes(getUserTable(user)).AutoMigrate(&domain.User{}).Create(user)
	return nil
}
