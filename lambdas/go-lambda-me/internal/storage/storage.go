package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go-lambda-me/pkg/domain"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	db.AutoMigrate(&domain.Payment{})
	return &PaymentRepository{
		db: db,
	}
}
func getPaymentTable() func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// TODO: Must receive the user to extract company name to specify the table name
		return tx.Table("payments")
	}
}
func getUsersTable() func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// TODO: Must receive the user to extract company name to specify the table name
		return tx.Table("users")
	}
}

func (r *PaymentRepository) GetAllPayments(userId string) ([]domain.Payment, error) {
	var payments []domain.Payment
	query := r.db.Scopes(getPaymentTable()).Where("user_id = ?", userId).Order("created_date")

	// Execute the query
	err := query.Find(&payments).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("No payments found with id: ", userId)
	}
	if err != nil {
		fmt.Println("Error getting payments:", err)
		return nil, err
	}

	return payments, nil
}

func (r *PaymentRepository) GetUserInformation(email string) (domain.User, error) {
	var user domain.User
	query := r.db.Scopes(getUsersTable()).Where("email = ?", email)
	err := query.First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("No user found with email: ", email)
		return user, err
	}
	return user, nil
}
