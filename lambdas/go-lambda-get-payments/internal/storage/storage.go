package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go-lambda-get-payments/pkg/domain"
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

func (r *PaymentRepository) GetAllPayments(filters map[string]string) ([]domain.Payment, error) {
	var payments []domain.Payment
	query := r.db.Scopes(getPaymentTable())

	// TODO: Always filter by confirmed users
	// Apply filters to the query
	if paymentType, ok := filters["payment_type"]; ok {
		query = query.Where("Type = ?", paymentType)
	}
	if hasReceipt, ok := filters["receipt"]; ok {
		if hasReceipt == "true" {
			query = query.Where("receipt <> ''")
		} else if hasReceipt == "false" {
			query = query.Where("receipt = '' OR receipt IS NULL")
		}
	}

	query = query.Order("created_date")

	// Execute the query
	err := query.Find(&payments).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("No payments found")
		return nil, nil
	}
	if err != nil {
		fmt.Println("Error getting payments:", err)
		return nil, err
	}

	return payments, nil
}
