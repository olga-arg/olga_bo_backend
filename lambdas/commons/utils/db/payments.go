package db

import (
	"commons/domain"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type PaymentRepository struct {
	Db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	db.AutoMigrate(&domain.Payment{})
	return &PaymentRepository{
		Db: db,
	}
}

func getPaymentTable(companyID string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		tableName := fmt.Sprintf("%s_payments", companyID)
		return tx.Table(tableName)
	}
}

func (r *PaymentRepository) Save(payment *domain.Payment, companyId string) error {
	err := r.Db.Scopes(getPaymentTable(companyId)).Create(payment).Error
	if err != nil {
		fmt.Println("Error saving Payment: ", err)
		return err
	}
	return nil
}

func (r *PaymentRepository) GetAllPayments(filters map[string]string, companyId string) ([]domain.Payment, error) {
	var payments domain.Payments
	query := r.Db.Scopes(getPaymentTable(companyId)).Preload("User")

	// TODO: Always filter by confirmed users
	// Apply filters to the query
	if paymentType, ok := filters["payment_type"]; ok {
		query = query.Where("Type = ?", paymentType)
	}
	if hasReceipt, ok := filters["receipt"]; ok {
		if hasReceipt == "true" {
			query = query.Where("receipt_image_key <> ''")
		} else if hasReceipt == "false" {
			query = query.Where("receipt_image_key = '' OR receipt IS NULL")
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
