package db

import (
	"commons/domain"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
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

func (r *UserRepository) Save(user *domain.User, companyId string) error {
	err := r.Db.Scopes(getUserTable(companyId)).Create(user).Error
	if err != nil {
		fmt.Println("Error saving user: ", err)
		return err
	}
	return nil
}

func (r *TeamRepository) GetReviewerById(id, companyId string) error {
	var user domain.User
	err := r.Db.Scopes(getUserTable(companyId)).Where("id = ?", id).First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return fmt.Errorf("user not found: %s", id)
		}
		fmt.Println("Error getting user by id: ", err)
		return err
	}
	return nil
}

func (r *UserRepository) GetAllUsers(filters map[string]string, companyId string) ([]domain.User, error) {
	var users []domain.User
	query := r.Db.Scopes(getUserTable(companyId)).
		Preload("Teams", func(db *gorm.DB) *gorm.DB {
			return db.Scopes(getTeamTable(companyId))
		})

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
		fmt.Println("No users found")
		return nil, nil
	}
	if err != nil {
		fmt.Println("Error getting users:", err)
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) UpdateUser(newUser *domain.User, companyId string) error {
	// Save the updated user
	err := r.Db.Scopes(getUserTable(companyId)).Save(newUser).Error
	if err != nil {
		fmt.Println("Error updating user:", err)
		return errors.Wrap(err, "failed to update user")
	}
	return nil
}

func (r *UserRepository) GetUserByID(userID, companyId string) (*domain.User, error) {
	var user domain.User
	query := r.Db.Scopes(getUserTable(companyId)).Where("id = ?", userID)
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
